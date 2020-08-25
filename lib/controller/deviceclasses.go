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
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Controller) ReadDeviceClass(jwt jwt_http_router.Jwt, id string) (deviceClass model.DeviceClass, err error, code int) {
	return this.com.GetDeviceClass(jwt, id)
}

func (this *Controller) PublishDeviceClassCreate(jwt jwt_http_router.Jwt, deviceClass model.DeviceClass) (model.DeviceClass, error, int) {
	if !com.IsAdmin(jwt) {
		return deviceClass, errors.New("access denied"), http.StatusForbidden
	}
	deviceClass.GenerateId()
	err, code := this.com.ValidateDeviceClass(jwt, deviceClass)
	if err != nil {
		return deviceClass, err, code
	}
	err = this.publisher.PublishDeviceClass(deviceClass, jwt.UserId)
	if err != nil {
		return deviceClass, err, http.StatusInternalServerError
	}
	return deviceClass, nil, http.StatusOK
}

func (this *Controller) PublishDeviceClassUpdate(jwt jwt_http_router.Jwt, id string, deviceClass model.DeviceClass) (model.DeviceClass, error, int) {
	if !com.IsAdmin(jwt) {
		return deviceClass, errors.New("access denied"), http.StatusForbidden
	}
	if deviceClass.Id != id {
		return deviceClass, errors.New("device id in body unequal to device id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	deviceClass.GenerateId()
	deviceClass.Id = id

	err, code := this.com.ValidateDeviceClass(jwt, deviceClass)
	if err != nil {
		return deviceClass, err, code
	}
	err = this.publisher.PublishDeviceClass(deviceClass, jwt.UserId)
	if err != nil {
		return deviceClass, err, http.StatusInternalServerError
	}
	return deviceClass, nil, http.StatusOK
}

func (this *Controller) PublishDeviceClassDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	if !com.IsAdmin(jwt) {
		return errors.New("access denied"), http.StatusForbidden
	}
	err := this.publisher.PublishDeviceClassDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
