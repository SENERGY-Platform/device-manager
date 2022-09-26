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

package com

import (
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"net/http"
)

func (this *Com) GetConcept(token auth.Token, id string) (concept model.Concept, err error, code int) {
	err, code = getResourceFromService(token, this.config.DeviceRepoUrl+"/concepts", id, &concept)
	return
}

func (this *Com) ValidateConcept(token auth.Token, concept model.Concept) (err error, code int) {
	if err = PreventIdModifier(concept.Id); err != nil {
		return err, http.StatusBadRequest
	}
	err, code = validateResources(token, this.config, []string{this.config.DeviceRepoUrl + "/concepts?dry-run=true"}, concept)
	if err != nil {
		return err, code
	}
	if this.config.ConverterUrl != "" && this.config.ConverterUrl != "-" {
		err, code = validateResource(token, this.config, "POST", this.config.ConverterUrl+"/validate/extended-conversions", map[string]interface{}{
			"nodes":      concept.CharacteristicIds,
			"extensions": concept.Conversions,
		})
	}
	return err, code
}

func (this *Com) ValidateConceptDelete(token auth.Token, id string) (err error, code int) {
	if err = PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	return validateResourceDelete(token, this.config, []string{
		this.config.DeviceRepoUrl + "/concepts",
	}, id)
}
