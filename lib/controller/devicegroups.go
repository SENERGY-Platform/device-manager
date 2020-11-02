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

package controller

import (
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"runtime/debug"
)

func (this *Controller) ReadDeviceGroup(jwt jwt_http_router.Jwt, id string) (dt model.DeviceGroup, err error, code int) {
	return this.com.GetTechnicalDeviceGroup(jwt, id)
}

func (this *Controller) PublishDeviceGroupCreate(jwt jwt_http_router.Jwt, dt model.DeviceGroup) (model.DeviceGroup, error, int) {
	dt.GenerateId()
	err, code := this.com.ValidateDeviceGroup(jwt, dt)
	if err != nil {
		return dt, err, code
	}
	err = this.publisher.PublishDeviceGroup(dt, jwt.UserId)
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceGroupUpdate(jwt jwt_http_router.Jwt, id string, dt model.DeviceGroup) (model.DeviceGroup, error, int) {
	if dt.Id != id {
		return dt, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	dt.GenerateId()

	if !com.IsAdmin(jwt) {
		err, code := this.com.PermissionCheckForDeviceGroup(jwt, id, "w")
		if err != nil {
			debug.PrintStack()
			return dt, err, code
		}
	}
	err, code := this.com.ValidateDeviceGroup(jwt, dt)
	if err != nil {
		debug.PrintStack()
		return dt, err, code
	}
	err = this.publisher.PublishDeviceGroup(dt, jwt.UserId)
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceGroupDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	exists, err, code := this.com.DevicesOfTypeExist(jwt, id)
	if err != nil {
		return err, code
	}
	if exists {
		return errors.New("expect no dependent devices"), http.StatusBadRequest
	}
	err, code = this.com.PermissionCheckForDeviceGroup(jwt, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceGroupDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
