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
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Controller) DeviceLocalIdToId(jwt jwt_http_router.Jwt, localId string) (id string, err error, errCode int) {
	return this.com.DeviceLocalIdToId(jwt, localId)
}

func (this *Controller) ReadDevice(jwt jwt_http_router.Jwt, id string) (device model.Device, err error, code int) {
	return this.com.GetDevice(jwt, id)
}

func (this *Controller) PublishDeviceCreate(jwt jwt_http_router.Jwt, device model.Device) (model.Device, error, int) {
	device.GenerateId()
	err, code := this.com.ValidateDevice(jwt, device)
	if err != nil {
		return device, err, code
	}
	err = this.publisher.PublishDevice(device, jwt.UserId)
	if err != nil {
		return device, err, http.StatusInternalServerError
	}
	return device, nil, http.StatusOK
}

func (this *Controller) PublishDeviceUpdate(jwt jwt_http_router.Jwt, id string, device model.Device) (model.Device, error, int) {
	if device.Id != id {
		return device, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	device.GenerateId()
	device.Id = id

	err, code := this.com.PermissionCheckForDevice(jwt, id, "w")
	if err != nil {
		return device, err, code
	}
	err, code = this.com.ValidateDevice(jwt, device)
	if err != nil {
		return device, err, code
	}
	err = this.publisher.PublishDevice(device, jwt.UserId)
	if err != nil {
		return device, err, http.StatusInternalServerError
	}
	return device, nil, http.StatusOK
}

func (this *Controller) PublishDeviceDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	err, code := this.com.PermissionCheckForDevice(jwt, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
