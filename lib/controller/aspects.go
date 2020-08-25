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

func (this *Controller) ReadAspect(jwt jwt_http_router.Jwt, id string) (aspect model.Aspect, err error, code int) {
	return this.com.GetAspect(jwt, id)
}

func (this *Controller) PublishAspectCreate(jwt jwt_http_router.Jwt, aspect model.Aspect) (model.Aspect, error, int) {
	if !com.IsAdmin(jwt) {
		return aspect, errors.New("access denied"), http.StatusForbidden
	}
	aspect.GenerateId()
	err, code := this.com.ValidateAspect(jwt, aspect)
	if err != nil {
		return aspect, err, code
	}
	err = this.publisher.PublishAspect(aspect, jwt.UserId)
	if err != nil {
		return aspect, err, http.StatusInternalServerError
	}
	return aspect, nil, http.StatusOK
}

func (this *Controller) PublishAspectUpdate(jwt jwt_http_router.Jwt, id string, aspect model.Aspect) (model.Aspect, error, int) {
	if !com.IsAdmin(jwt) {
		return aspect, errors.New("access denied"), http.StatusForbidden
	}
	if aspect.Id != id {
		return aspect, errors.New("device id in body unequal to device id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	aspect.GenerateId()
	aspect.Id = id

	err, code := this.com.ValidateAspect(jwt, aspect)
	if err != nil {
		return aspect, err, code
	}
	err = this.publisher.PublishAspect(aspect, jwt.UserId)
	if err != nil {
		return aspect, err, http.StatusInternalServerError
	}
	return aspect, nil, http.StatusOK
}

func (this *Controller) PublishAspectDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	if !com.IsAdmin(jwt) {
		return errors.New("access denied"), http.StatusForbidden
	}
	err := this.publisher.PublishAspectDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
