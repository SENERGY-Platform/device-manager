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
	"runtime/debug"
	"sort"
)

func (this *Controller) ReadDeviceType(token auth.Token, id string) (dt model.DeviceType, err error, code int) {
	dt, err, code = this.com.GetDeviceType(token, id)
	sort.Slice(dt.Services, func(i, j int) bool {
		return dt.Services[i].Name < dt.Services[j].Name
	})
	return dt, err, code
}

func (this *Controller) PublishDeviceTypeCreate(token auth.Token, dt model.DeviceType) (model.DeviceType, error, int) {
	dt.GenerateId()
	err, code := this.com.ValidateDeviceType(token, dt)
	if err != nil {
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, token.GetUserId())
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeUpdate(token auth.Token, id string, dt model.DeviceType) (model.DeviceType, error, int) {
	if dt.Id != id {
		return dt, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	dt.GenerateId()

	if !token.IsAdmin() {
		err, code := this.com.PermissionCheckForDeviceType(token, id, "w")
		if err != nil {
			debug.PrintStack()
			return dt, err, code
		}
	}
	err, code := this.com.ValidateDeviceType(token, dt)
	if err != nil {
		debug.PrintStack()
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, token.GetUserId())
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeDelete(token auth.Token, id string) (error, int) {
	exists, err, code := this.com.DevicesOfTypeExist(token, id)
	if err != nil {
		return err, code
	}
	if exists {
		return errors.New("expect no dependent devices"), http.StatusBadRequest
	}
	err, code = this.com.PermissionCheckForDeviceType(token, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceTypeDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
