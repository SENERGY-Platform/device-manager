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

// expects previous permission check and use own admin jwt to access device
func (this *Com) GetDevice(token auth.Token, id string) (device models.Device, err error, code int) {
	err, code = getResourceFromService(token, this.config.DeviceRepoUrl+"/devices", id, &device)
	return
}

func (this *Com) ValidateDevice(token auth.Token, device models.Device) (err error, code int) {
	if err = PreventIdModifier(device.Id); err != nil {
		return err, http.StatusBadRequest
	}
	return validateResources(token, this.config, []string{
		this.config.DeviceRepoUrl + "/devices?dry-run=true",
	}, device)
}
