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
	"fmt"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
)

func (this *Controller) ListDevicesByQuery(token auth.Token, query url.Values) (devices []models.Device, err error, code int) {
	return this.com.ListDevicesByQuery(token, query)
}

func (this *Controller) DeviceLocalIdToId(token auth.Token, ownerId string, localId string) (id string, err error, errCode int) {
	device, err, code := this.com.GetDeviceByLocalId(token, ownerId, localId)
	return device.Id, err, code
}

func (this *Controller) ReadDeviceByLocalId(token auth.Token, ownerId string, localId string) (device models.Device, err error, errCode int) {
	return this.com.GetDeviceByLocalId(token, ownerId, localId)
}

func (this *Controller) ReadDevice(token auth.Token, id string) (device models.Device, err error, code int) {
	return this.com.GetDevice(token, id)
}

func (this *Controller) PublishDeviceCreate(token auth.Token, device models.Device, options model.DeviceCreateOptions) (models.Device, error, int) {
	device.GenerateId()
	if device.OwnerId != "" && device.OwnerId != token.GetUserId() {
		return device, errors.New("new devices must be initialised with the requesting user as owner-id"), http.StatusBadRequest
	}
	device.OwnerId = token.GetUserId()

	err, code := this.com.ValidateDevice(token, device)
	if err != nil {
		return device, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceTopic,
		ResourceId:   device.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishDevice(device, token.GetUserId())
	if err != nil {
		return device, err, http.StatusInternalServerError
	}

	err = wait()
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

	var original models.Device
	var exists bool
	original, err, code = this.com.GetDevice(token, device.Id)
	if err != nil && code != http.StatusNotFound {
		return device, err, code
	}
	if err != nil {
		err, code = nil, 200
		exists = false
	} else {
		exists = true
	}

	if exists && len(options.UpdateOnlySameOriginAttributes) > 0 {
		device.Attributes = updateSameOriginAttributes(original.Attributes, device.Attributes, options.UpdateOnlySameOriginAttributes)
	}

	//set device owner-id if none is given
	//prefer existing owner, fallback to requesting user
	if device.OwnerId == "" {
		device.OwnerId = original.OwnerId //may be empty for new devices
	}
	if device.OwnerId == "" {
		device.OwnerId = token.GetUserId()
	}

	if exists && original.OwnerId != device.OwnerId && original.OwnerId != "" && !token.IsAdmin() {
		err, code := this.com.PermissionCheckForDevice(token, device.Id, "a")
		if err != nil {
			if code == http.StatusForbidden {
				return device, fmt.Errorf("only admins may set new owner: %w", err), http.StatusBadRequest
			} else {
				return device, err, code
			}
		}
	}

	err, code = this.com.ValidateDevice(token, device)
	if err != nil {
		return device, err, code
	}

	rights, found, err := this.com.GetResourceRights(token, this.config.DeviceTopic, device.Id, "w")
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return device, err, http.StatusInternalServerError
	}

	//new device owner-id must be existing admin user (ignore for new devices or devices with unchanged owner)
	if found && device.OwnerId != original.OwnerId && !slices.Contains(rights.PermissionHolders.AdminUsers, device.OwnerId) {
		return device, errors.New("new owner must have existing user admin rights"), http.StatusBadRequest
	}

	//ensure retention of original creator
	creator := rights.Creator
	if !found || creator == "" {
		creator = token.GetUserId()
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceTopic,
		ResourceId:   device.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishDevice(device, creator)
	if err != nil {
		return device, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return device, err, http.StatusInternalServerError
	}

	return device, nil, http.StatusOK
}

func (this *Controller) PublishDeviceDelete(token auth.Token, id string, options model.DeviceDeleteOptions) (error, int) {
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	err, code := this.com.PermissionCheckForDevice(token, id, "a")
	if err != nil {
		return err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishDeviceDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
