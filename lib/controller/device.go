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
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"net/http"
)

func (this *Controller) DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, errCode int) {
	return this.com.DeviceLocalIdToId(token, localId)
}

func (this *Controller) ReadDevice(token auth.Token, id string) (device model.Device, err error, code int) {
	return this.com.GetDevice(token, id)
}

func (this *Controller) PublishDeviceCreate(token auth.Token, device model.Device) (model.Device, error, int) {
	device.GenerateId()
	err, code := this.com.ValidateDevice(token, device)
	if err != nil {
		return device, err, code
	}
	err = this.publisher.PublishDevice(device, token.GetUserId())
	if err != nil {
		return device, err, http.StatusInternalServerError
	}
	return device, nil, http.StatusOK
}

func (this *Controller) PublishDeviceUpdate(token auth.Token, id string, device model.Device) (model.Device, error, int) {
	if device.Id != id {
		return device, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	device.GenerateId()
	device.Id = id

	if !token.IsAdmin() {
		err, code := this.com.PermissionCheckForDevice(token, id, "w")
		if err != nil {
			return device, err, code
		}
	}
	err, code := this.com.ValidateDevice(token, device)
	if err != nil {
		return device, err, code
	}
	err = this.publisher.PublishDevice(device, token.GetUserId())
	if err != nil {
		return device, err, http.StatusInternalServerError
	}
	return device, nil, http.StatusOK
}

func (this *Controller) PublishDeviceDelete(token auth.Token, id string) (error, int) {
	err, code := this.com.PermissionCheckForDevice(token, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
