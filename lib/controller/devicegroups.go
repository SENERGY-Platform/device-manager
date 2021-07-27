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
	"net/http"
	"runtime/debug"
)

func (this *Controller) ReadDeviceGroup(token string, id string) (dt model.DeviceGroup, err error, code int) {
	return this.com.GetTechnicalDeviceGroup(token, id)
}

func (this *Controller) PublishDeviceGroupCreate(token string, dg model.DeviceGroup) (model.DeviceGroup, error, int) {
	dg.GenerateId()
	dg.SetShortCriteria()

	err, code := this.com.ValidateDeviceGroup(token, dg)
	if err != nil {
		return dg, err, code
	}
	err = this.publisher.PublishDeviceGroup(dg, com.GetUserId(token))
	if err != nil {
		return dg, err, http.StatusInternalServerError
	}
	return dg, nil, http.StatusOK
}

func (this *Controller) PublishDeviceGroupUpdate(token string, id string, dg model.DeviceGroup) (model.DeviceGroup, error, int) {
	if dg.Id != id {
		return dg, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	dg.GenerateId()
	dg.SetShortCriteria()

	if !com.IsAdmin(token) {
		err, code := this.com.PermissionCheckForDeviceGroup(token, id, "w")
		if err != nil {
			debug.PrintStack()
			return dg, err, code
		}
	}
	err, code := this.com.ValidateDeviceGroup(token, dg)
	if err != nil {
		debug.PrintStack()
		return dg, err, code
	}
	err = this.publisher.PublishDeviceGroup(dg, com.GetUserId(token))
	if err != nil {
		debug.PrintStack()
		return dg, err, http.StatusInternalServerError
	}
	return dg, nil, http.StatusOK
}

func (this *Controller) PublishDeviceGroupDelete(token string, id string) (error, int) {
	err, code := this.com.PermissionCheckForDeviceGroup(token, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceGroupDelete(id, com.GetUserId(token))
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
