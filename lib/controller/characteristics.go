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

func (this *Controller) PublishCharacteristicUpdate(jwt jwt_http_router.Jwt, id string, concept model.Concept) (model.Concept, error, int) {
	if concept.Id != id {
		return concept, errors.New("concept id in body unequal to concept id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	concept.GenerateId()

	err, code := this.com.PermissionCheckForDeviceType(jwt, id, "w")
	if err != nil {
		debug.PrintStack()
		return concept, err, code
	}
	err, code = this.com.ValidateConcept(jwt, concept)
	if err != nil {
		debug.PrintStack()
		return concept, err, code
	}
	err = this.publisher.PublishConcept(concept, jwt.UserId)
	if err != nil {
		debug.PrintStack()
		return concept, err, http.StatusInternalServerError
	}
	return concept, nil, http.StatusOK
}

func (this *Controller) PublishCharacteristicDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	err, code := this.com.PermissionCheckForConcept(jwt, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishConceptDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
