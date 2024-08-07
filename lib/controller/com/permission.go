/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/permission-search/lib/client"
	permmodel "github.com/SENERGY-Platform/permission-search/lib/model"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *Com) GetPermissions(token auth.Token, kind string, id string) (permmodel.ResourceRights, error) {
	return this.perm.GetRights(token.Jwt(), kind, id)
}

// GetResourceOwner queries the permission-search service for the entity identified by kind and id and extracts its owner
// the rights parameter is a mandatory part of the permission-search api
// it is used to identify which rights the user (token) must have for the entity, to get the entity as a result
// for example, if a user has 'r' rights to an entity, the query will find the entity, if requested with rights="r" but not with rights="w" or rights="rw"
func (this *Com) GetResourceOwner(token auth.Token, kind string, id string, rights string) (owner string, found bool, err error) {
	temp, _, err := client.Query[[]permmodel.EntryResult](this.perm, token.Jwt(), client.QueryMessage{
		Resource: kind,
		ListIds: &permmodel.QueryListIds{
			QueryListCommons: permmodel.QueryListCommons{
				Limit:  1,
				Offset: 0,
				Rights: rights,
			},
			Ids: []string{id},
		},
	})
	if err != nil {
		return owner, false, err
	}
	if len(temp) == 0 {
		return owner, false, nil
	}
	return temp[0].Creator, true, nil
}

func (this *Com) GetResourceRights(token auth.Token, kind string, id string, rights string) (result permmodel.EntryResult, found bool, err error) {
	temp, _, err := client.Query[[]permmodel.EntryResult](this.perm, token.Jwt(), client.QueryMessage{
		Resource: kind,
		ListIds: &permmodel.QueryListIds{
			QueryListCommons: permmodel.QueryListCommons{
				Limit:  1,
				Offset: 0,
				Rights: rights,
			},
			Ids: []string{id},
		},
	})
	if err != nil {
		return result, false, err
	}
	if len(temp) == 0 {
		return result, false, nil
	}
	return temp[0], true, nil
}

func (this *Com) PermissionCheckForDeviceList(token auth.Token, ids []string, rights string) (result map[string]bool, err error, code int) {
	ids = append(ids, removeIdModifiers(ids)...)
	ids = RemoveDuplicates(ids)
	result, code, err = client.Query[map[string]bool](this.perm, token.Jwt(), client.QueryMessage{
		Resource: "devices",
		CheckIds: &client.QueryCheckIds{
			Ids:    ids,
			Rights: rights,
		},
	})
	return
}

func (this *Com) PermissionCheckForDevice(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.DeviceTopic)
}

func (this *Com) PermissionCheckForHub(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.HubTopic)
}

func (this *Com) PermissionCheckForDeviceGroup(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.DeviceGroupTopic)
}

func (this *Com) PermissionCheckForDeviceType(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.DeviceTypeTopic)
}

func (this *Com) PermissionCheckForConcept(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.ConceptTopic)
}

func (this *Com) PermissionCheckForCharacteristic(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.CharacteristicTopic)
}

func (this *Com) PermissionCheckForLocation(token auth.Token, id string, permission string) (err error, code int) {
	if token.IsAdmin() {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(token, id, permission, this.config.LocationTopic)
}

func (this *Com) PermissionCheck(token auth.Token, id string, permission string, resource string) (err error, code int) {
	if this.config.PermissionsUrl == "" || this.config.PermissionsUrl == "-" {
		return nil, 200
	}
	id = removeIdModifier(id)
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/v3/resources/"+url.QueryEscape(resource)+"/"+url.QueryEscape(id)+"/access?rights="+url.QueryEscape(permission), nil)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		resp.Body.Close()
		log.Println("DEBUG: PermissionCheck()", buf.String())
		err = errors.New("access denied")
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}

	var ok bool
	err = json.NewDecoder(resp.Body).Decode(&ok)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	if !ok {
		return errors.New("access denied"), http.StatusForbidden
	}
	return
}

func (this *Com) DevicesOfTypeExist(token auth.Token, deviceTypeId string) (result bool, err error, code int) {
	if !token.IsAdmin() {
		return false, errors.New("only for admins allowed"), http.StatusForbidden
	}
	if this.config.PermissionsUrl == "" || this.config.PermissionsUrl == "-" {
		return false, nil, 200
	}
	deviceTypeId = removeIdModifier(deviceTypeId)
	endpoint := this.config.PermissionsUrl + "/v3/resources/devices?limit=1&rights=x&filter=" + url.QueryEscape("device_type_id:"+deviceTypeId)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		debug.PrintStack()
		return result, errors.New(buf.String()), resp.StatusCode
	}
	temp := []interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	return len(temp) > 0, nil, http.StatusOK
}

func (this *Com) DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, code int) {
	endpoint := this.config.PermissionsUrl + "/v3/resources/devices?limit=1&rights=x&filter=" + url.QueryEscape("local_id:"+localId)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		debug.PrintStack()
		return "", err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return "", err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		debug.PrintStack()
		return "", errors.New(buf.String()), resp.StatusCode
	}
	temp := []map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		debug.PrintStack()
		return "", err, http.StatusInternalServerError
	}
	if len(temp) == 0 {
		return "", errors.New("not found"), http.StatusNotFound
	}
	if idinterface, ok := temp[0]["id"]; !ok {
		err = errors.New("id field not found")
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", err, http.StatusInternalServerError
	} else if id, ok = idinterface.(string); !ok {
		err = errors.New("id field is not string")
		log.Println("ERROR:", err)
		debug.PrintStack()
		return "", err, http.StatusInternalServerError
	} else {
		return id, nil, http.StatusOK
	}
}
