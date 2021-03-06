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
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

func IsAdmin(jwt jwt_http_router.Jwt) bool {
	return contains(jwt.RealmAccess.Roles, "admin")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (this *Com) PermissionCheckForDevice(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.DeviceTopic)
}

func (this *Com) PermissionCheckForHub(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.HubTopic)
}

func (this *Com) PermissionCheckForDeviceGroup(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.DeviceGroupTopic)
}

func (this *Com) PermissionCheckForDeviceType(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.DeviceTypeTopic)
}

func (this *Com) PermissionCheckForConcept(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.ConceptTopic)
}

func (this *Com) PermissionCheckForCharacteristic(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.CharacteristicTopic)
}

func (this *Com) PermissionCheckForLocation(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	if IsAdmin(jwt) {
		return nil, http.StatusOK
	}
	return this.PermissionCheck(jwt, id, permission, this.config.LocationTopic)
}

func (this *Com) PermissionCheck(jwt jwt_http_router.Jwt, id string, permission string, resource string) (err error, code int) {
	if this.config.PermissionsUrl == "" || this.config.PermissionsUrl == "-" {
		return nil, 200
	}
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/jwt/check/"+url.QueryEscape(resource)+"/"+url.QueryEscape(id)+"/"+permission+"/bool", nil)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", string(jwt.Impersonate))
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

func (this *Com) DevicesOfTypeExist(jwt jwt_http_router.Jwt, deviceTypeId string) (result bool, err error, code int) {
	if !IsAdmin(jwt) {
		return false, errors.New("only for admins allowed"), http.StatusForbidden
	}
	if this.config.PermissionsUrl == "" || this.config.PermissionsUrl == "-" {
		return false, nil, 200
	}
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/jwt/select/devices/device_type_id/"+url.PathEscape(deviceTypeId)+"/x", nil)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", string(jwt.Impersonate))
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

func (this *Com) DeviceLocalIdToId(jwt jwt_http_router.Jwt, localId string) (id string, err error, code int) {
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/jwt/select/devices/local_id/"+url.PathEscape(localId)+"/x", nil)
	if err != nil {
		debug.PrintStack()
		return "", err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", string(jwt.Impersonate))
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
