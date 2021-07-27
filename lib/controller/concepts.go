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
	"net/http"
	"runtime/debug"
)

func (this *Controller) PublishConceptCreate(token string, concept model.Concept) (model.Concept, error, int) {
	concept.GenerateId()
	err, code := this.com.ValidateConcept(token, concept)
	if err != nil {
		return concept, err, code
	}
	err = this.publisher.PublishConcept(concept, com.GetUserId(token))
	if err != nil {
		return concept, err, http.StatusInternalServerError
	}
	return concept, nil, http.StatusOK
}

func (this *Controller) PublishConceptUpdate(token string, id string, concept model.Concept) (model.Concept, error, int) {
	if concept.Id != id {
		return concept, errors.New("concept id in body unequal to concept id in request endpoint"), http.StatusBadRequest
	}

	concept.GenerateId()
	if !com.IsAdmin(token) {
		err, code := this.com.PermissionCheckForConcept(token, id, "w")
		if err != nil {
			debug.PrintStack()
			return concept, err, code
		}
	}

	err, code := this.com.ValidateConcept(token, concept)
	if err != nil {
		debug.PrintStack()
		return concept, err, code
	}
	err = this.publisher.PublishConcept(concept, com.GetUserId(token))
	if err != nil {
		debug.PrintStack()
		return concept, err, http.StatusInternalServerError
	}
	return concept, nil, http.StatusOK
}

func (this *Controller) PublishConceptDelete(token string, id string) (error, int) {
	err, code := this.com.PermissionCheckForConcept(token, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishConceptDelete(id, com.GetUserId(token))
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
