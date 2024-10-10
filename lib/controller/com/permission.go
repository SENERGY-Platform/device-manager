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
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	permsearch "github.com/SENERGY-Platform/permission-search/lib/client"
	permmodel "github.com/SENERGY-Platform/permission-search/lib/model"
	"github.com/SENERGY-Platform/permissions-v2/pkg/model"
	"log"
	"net/http"
)

// GetResourceOwner queries the permission-search service for the entity identified by kind and id and extracts its owner
// the rights parameter is a mandatory part of the permission-search api
// it is used to identify which rights the user (token) must have for the entity, to get the entity as a result
// for example, if a user has 'r' rights to an entity, the query will find the entity, if requested with rights="r" but not with rights="w" or rights="rw"
// TODO: replace
func (this *Com) GetResourceOwner(token auth.Token, kind string, id string, rights string) (owner string, found bool, err error) {
	temp, _, err := permsearch.Query[[]permmodel.EntryResult](this.search, token.Jwt(), permsearch.QueryMessage{
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

// TODO: replace
func (this *Com) GetResourceRights(token auth.Token, kind string, id string, rights string) (result permmodel.EntryResult, found bool, err error) {
	temp, _, err := permsearch.Query[[]permmodel.EntryResult](this.search, token.Jwt(), permsearch.QueryMessage{
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
	permissions, err := model.PermissionListFromString(rights)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return this.perm.CheckMultiplePermissions(token.Jwt(), this.config.DeviceTopic, ids, permissions...)
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
	permissions, err := model.PermissionListFromString(permission)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	log.Printf("DEBUG: PermissionCheck %#v %#v %#v", resource, id, permissions) //TODO: remove
	access, err, code := this.perm.CheckPermission(token.Jwt(), resource, id, permissions...)
	if err != nil {
		return err, code
	}
	if !access {
		return errors.New("access denied"), http.StatusForbidden
	}
	return nil, http.StatusOK
}

func (this *Com) DevicesOfTypeExist(token auth.Token, deviceTypeId string) (result bool, err error, code int) {
	if !token.IsAdmin() {
		return false, errors.New("only for admins allowed"), http.StatusForbidden
	}
	deviceTypeId = removeIdModifier(deviceTypeId)
	devices, err, code := this.devices.ListDevices(token.Jwt(), devicerepo.DeviceListOptions{
		DeviceTypeIds: []string{deviceTypeId},
		Limit:         1,
		Offset:        0,
	})
	if err != nil {
		return false, err, code
	}
	return len(devices) > 0, nil, http.StatusOK
}

func (this *Com) DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, code int) {
	device, err, code := this.devices.ReadDeviceByLocalId(token.GetUserId(), localId, token.Jwt(), devicerepo.READ)
	if err != nil {
		return "", err, code
	}
	return device.Id, nil, http.StatusOK
}
