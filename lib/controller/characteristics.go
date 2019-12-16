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
)

func (this *Controller) PublishCharacteristicCreate(jwt jwt_http_router.Jwt, conceptId string, characteristic model.Characteristic) (model.Characteristic, error, int) {
	characteristic.GenerateId()
	err, code := this.com.ValidateCharacteristic(jwt, characteristic)
	if err != nil {
		return characteristic, err, code
	}
	err = this.publisher.PublishCharacteristic(conceptId, characteristic, jwt.UserId)
	if err != nil {
		return characteristic, err, http.StatusInternalServerError
	}
	return characteristic, nil, http.StatusOK
}

func (this *Controller) PublishCharacteristicUpdate(jwt jwt_http_router.Jwt, conceptId string, characteristicId string, characteristic model.Characteristic) (model.Characteristic, error, int) {
	if characteristic.Id != characteristicId {
		return characteristic, errors.New("characteristic id in body unequal to characteristic id in request endpoint"), http.StatusBadRequest
	}

	characteristic.GenerateId()

	if !com.IsAdmin(jwt){
		err, code := this.com.PermissionCheckForCharacteristic(jwt, characteristicId, "w")
		if err != nil {
			debug.PrintStack()
			return characteristic, err, code
		}
	}
	err, code := this.com.ValidateCharacteristic(jwt, characteristic)
	if err != nil {
		debug.PrintStack()
		return characteristic, err, code
	}
	err = this.publisher.PublishCharacteristic(conceptId, characteristic, jwt.UserId)
	if err != nil {
		debug.PrintStack()
		return characteristic, err, http.StatusInternalServerError
	}
	return characteristic, nil, http.StatusOK
}

func (this *Controller) PublishCharacteristicDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	err, code := this.com.PermissionCheckForCharacteristic(jwt, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishCharacteristicDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
