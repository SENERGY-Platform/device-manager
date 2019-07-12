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
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Controller) ReadHub(jwt jwt_http_router.Jwt, id string) (hub model.Hub, err error, code int) {
	return this.com.GetHub(jwt, id)
}

func (this *Controller) PublishHubCreate(jwt jwt_http_router.Jwt, hub model.Hub) (model.Hub, error, int) {
	hub.GenerateId()
	err, code := this.com.ValidateHub(jwt, hub)
	if err != nil {
		return hub, err, code
	}
	err = this.publisher.PublishHub(hub, jwt.UserId)
	if err != nil {
		return hub, err, http.StatusInternalServerError
	}
	return hub, nil, http.StatusOK
}

func (this *Controller) PublishHubUpdate(jwt jwt_http_router.Jwt, id string, hub model.Hub) (model.Hub, error, int) {
	if hub.Id != id {
		return hub, errors.New("hub id in body unequal to hub id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	hub.GenerateId()
	hub.Id = id

	err, code := this.com.PermissionCheckForHub(jwt, id, "w")
	if err != nil {
		return hub, err, code
	}
	err, code = this.com.ValidateHub(jwt, hub)
	if err != nil {
		return hub, err, code
	}
	err = this.publisher.PublishHub(hub, jwt.UserId)
	if err != nil {
		return hub, err, http.StatusInternalServerError
	}
	return hub, nil, http.StatusOK
}

func (this *Controller) PublishHubDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	err, code := this.com.PermissionCheckForHub(jwt, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishHubDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
