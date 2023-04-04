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
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"net/http"
)

func (this *Controller) DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, errCode int) {
	return this.com.DeviceLocalIdToId(token, localId)
}

func (this *Controller) ReadDevice(token auth.Token, id string) (device models.Device, err error, code int) {
	return this.com.GetDevice(token, id)
}

func (this *Controller) PublishDeviceCreate(token auth.Token, device models.Device) (models.Device, error, int) {
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

// admins may create new devices but only without setting options.UpdateOnlySameOriginAttributes
func (this *Controller) PublishDeviceUpdate(token auth.Token, id string, device models.Device, options model.DeviceUpdateOptions) (_ models.Device, err error, code int) {
	if device.Id != id {
		return device, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	if !token.IsAdmin() {
		err, code = this.com.PermissionCheckForDevice(token, id, "w")
		if err != nil {
			return device, err, code
		}
	}

	if len(options.UpdateOnlySameOriginAttributes) > 0 {
		var original models.Device
		original, err, code = this.com.GetDevice(token, device.Id)
		if err != nil {
			return device, err, code
		}
		device.Attributes = updateSameOriginAttributes(original.Attributes, device.Attributes, options.UpdateOnlySameOriginAttributes)
	}

	err, code = this.com.ValidateDevice(token, device)
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
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
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
