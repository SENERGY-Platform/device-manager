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
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Com) GetTechnicalDeviceGroup(jwt jwt_http_router.Jwt, id string) (dt model.DeviceGroup, err error, code int) {
	if this.config.DeviceRepoUrl == "" || this.config.DeviceRepoUrl == "-" {
		return model.DeviceGroup{}, nil, http.StatusOK
	}
	err, code = getResourceFromService(jwt, this.config.DeviceRepoUrl+"/device-groups", id, &dt)
	return
}

func (this *Com) ValidateDeviceGroup(jwt jwt_http_router.Jwt, dt model.DeviceGroup) (err error, code int) {
	list := []string{}
	if this.config.DeviceRepoUrl != "" && this.config.DeviceRepoUrl != "-" {
		list = append(list, this.config.DeviceRepoUrl+"/device-groups?dry-run=true")
	}
	return validateResource(jwt, list, dt)
}
