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
	"github.com/SENERGY-Platform/models/go/models"
	"net/http"
)

func (this *Com) GetTechnicalDeviceGroup(token auth.Token, id string) (dt models.DeviceGroup, err error, code int) {
	if this.config.DeviceRepoUrl == "" || this.config.DeviceRepoUrl == "-" {
		return models.DeviceGroup{}, nil, http.StatusOK
	}
	err, code = getResourceFromService(token, this.config.DeviceRepoUrl+"/device-groups", id, &dt)
	return
}

func (this *Com) ValidateDeviceGroup(token auth.Token, dg models.DeviceGroup) (err error, code int) {
	if err = PreventIdModifier(dg.Id); err != nil {
		return err, http.StatusBadRequest
	}
	list := []string{}
	if this.config.DeviceRepoUrl != "" && this.config.DeviceRepoUrl != "-" {
		list = append(list, this.config.DeviceRepoUrl+"/device-groups?dry-run=true")
	}
	return validateResources(token, this.config, list, dg)
}

func (this *Com) ValidateDeviceGroupDelete(token auth.Token, id string) (err error, code int) {
	if err = PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	return validateResourceDelete(token, this.config, []string{
		this.config.DeviceRepoUrl + "/device-groups",
	}, id)
}
