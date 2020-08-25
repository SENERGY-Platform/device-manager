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
	"sort"
)

func (this *Controller) ReadDeviceType(jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int) {
	tdt, err, code := this.com.GetTechnicalDeviceType(jwt, id)
	if err != nil {
		return tdt, err, code
	}
	sdt, err, code := this.com.GetSemanticDeviceType(jwt, id)
	if err != nil {
		return tdt, err, code
	}
	tdt.DeviceClassId = sdt.DeviceClassId
	index := map[string]model.Service{}
	for _, service := range sdt.Services {
		index[service.Id] = service
	}
	for i, service := range tdt.Services {
		service.FunctionIds = index[service.Id].FunctionIds
		service.AspectIds = index[service.Id].AspectIds
		tdt.Services[i] = service
	}

	sort.Slice(tdt.Services, func(i, j int) bool {
		return tdt.Services[i].Name < tdt.Services[j].Name
	})

	return tdt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeCreate(jwt jwt_http_router.Jwt, dt model.DeviceType) (model.DeviceType, error, int) {
	dt.GenerateId()
	err, code := this.com.ValidateDeviceType(jwt, dt)
	if err != nil {
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, jwt.UserId)
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeUpdate(jwt jwt_http_router.Jwt, id string, dt model.DeviceType) (model.DeviceType, error, int) {
	if dt.Id != id {
		return dt, errors.New("device id in body unequal to device id in request endpoint"), http.StatusBadRequest
	}

	dt.GenerateId()

	if !com.IsAdmin(jwt) {
		err, code := this.com.PermissionCheckForDeviceType(jwt, id, "w")
		if err != nil {
			debug.PrintStack()
			return dt, err, code
		}
	}
	err, code := this.com.ValidateDeviceType(jwt, dt)
	if err != nil {
		debug.PrintStack()
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, jwt.UserId)
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	exists, err, code := this.com.DevicesOfTypeExist(jwt, id)
	if err != nil {
		return err, code
	}
	if exists {
		return errors.New("expect no dependent devices"), http.StatusBadRequest
	}
	err, code = this.com.PermissionCheckForDeviceType(jwt, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceTypeDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
