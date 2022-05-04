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

func (this *Controller) ReadDeviceClass(token auth.Token, id string) (deviceClass model.DeviceClass, err error, code int) {
	return this.com.GetDeviceClass(token, id)
}

func (this *Controller) PublishDeviceClassCreate(token auth.Token, deviceClass model.DeviceClass) (model.DeviceClass, error, int) {
	if !token.IsAdmin() {
		return deviceClass, errors.New("access denied"), http.StatusForbidden
	}
	deviceClass.GenerateId()
	err, code := this.com.ValidateDeviceClass(token, deviceClass)
	if err != nil {
		return deviceClass, err, code
	}
	err = this.publisher.PublishDeviceClass(deviceClass, token.GetUserId())
	if err != nil {
		return deviceClass, err, http.StatusInternalServerError
	}
	return deviceClass, nil, http.StatusOK
}

func (this *Controller) PublishDeviceClassUpdate(token auth.Token, id string, deviceClass model.DeviceClass) (model.DeviceClass, error, int) {
	if !token.IsAdmin() {
		return deviceClass, errors.New("access denied"), http.StatusForbidden
	}
	if deviceClass.Id != id {
		return deviceClass, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	deviceClass.GenerateId()
	deviceClass.Id = id

	err, code := this.com.ValidateDeviceClass(token, deviceClass)
	if err != nil {
		return deviceClass, err, code
	}
	err = this.publisher.PublishDeviceClass(deviceClass, token.GetUserId())
	if err != nil {
		return deviceClass, err, http.StatusInternalServerError
	}
	return deviceClass, nil, http.StatusOK
}

func (this *Controller) PublishDeviceClassDelete(token auth.Token, id string) (error, int) {
	if !token.IsAdmin() {
		return errors.New("access denied"), http.StatusForbidden
	}
	err, code := this.com.ValidateDeviceClassDelete(token, id)
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceClassDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
