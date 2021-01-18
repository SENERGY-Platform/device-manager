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
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Controller) ReadLocation(jwt jwt_http_router.Jwt, id string) (Location model.Location, err error, code int) {
	return this.com.GetLocation(jwt, id)
}

func (this *Controller) PublishLocationCreate(jwt jwt_http_router.Jwt, Location model.Location) (model.Location, error, int) {
	if !com.IsAdmin(jwt) {
		return Location, errors.New("access denied"), http.StatusForbidden
	}
	Location.GenerateId()
	err, code := this.com.ValidateLocation(jwt, Location)
	if err != nil {
		return Location, err, code
	}
	err = this.publisher.PublishLocation(Location, jwt.UserId)
	if err != nil {
		return Location, err, http.StatusInternalServerError
	}
	return Location, nil, http.StatusOK
}

func (this *Controller) PublishLocationUpdate(jwt jwt_http_router.Jwt, id string, Location model.Location) (model.Location, error, int) {
	if !com.IsAdmin(jwt) {
		return Location, errors.New("access denied"), http.StatusForbidden
	}
	if Location.Id != id {
		return Location, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	Location.GenerateId()
	Location.Id = id

	err, code := this.com.ValidateLocation(jwt, Location)
	if err != nil {
		return Location, err, code
	}
	err = this.publisher.PublishLocation(Location, jwt.UserId)
	if err != nil {
		return Location, err, http.StatusInternalServerError
	}
	return Location, nil, http.StatusOK
}

func (this *Controller) PublishLocationDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	if !com.IsAdmin(jwt) {
		return errors.New("access denied"), http.StatusForbidden
	}
	err := this.publisher.PublishLocationDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
