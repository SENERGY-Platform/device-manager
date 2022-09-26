/*
 * Copyright 2021 InfAI (CC SES)
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
	"net/http"
)

func (this *Controller) ReadLocation(token auth.Token, id string) (location model.Location, err error, code int) {
	return this.com.GetLocation(token, id)
}

func (this *Controller) PublishLocationCreate(token auth.Token, location model.Location) (result model.Location, err error, code int) {
	location.GenerateId()
	location.DeviceIds, err = this.filterInvalidDeviceIds(token, location.DeviceIds, "r")
	if err != nil {
		return location, err, code
	}
	err, code = this.com.ValidateLocation(token, location)
	if err != nil {
		return location, err, code
	}
	err = this.publisher.PublishLocation(location, token.GetUserId())
	if err != nil {
		return location, err, http.StatusInternalServerError
	}
	return location, nil, http.StatusOK
}

func (this *Controller) PublishLocationUpdate(token auth.Token, id string, location model.Location) (model.Location, error, int) {
	if location.Id != id {
		return location, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	location.GenerateId()
	location.Id = id

	err, code := this.com.PermissionCheckForLocation(token, id, "w")
	if err != nil {
		return location, err, code
	}

	location.DeviceIds, err = this.filterInvalidDeviceIds(token, location.DeviceIds, "r")
	if err != nil {
		return location, err, code
	}

	err, code = this.com.ValidateLocation(token, location)
	if err != nil {
		return location, err, code
	}
	err = this.publisher.PublishLocation(location, token.GetUserId())
	if err != nil {
		return location, err, http.StatusInternalServerError
	}
	return location, nil, http.StatusOK
}

func (this *Controller) PublishLocationDelete(token auth.Token, id string) (error, int) {
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	err, code := this.com.PermissionCheckForLocation(token, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishLocationDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
