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
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/models/go/models"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *Com) ListDevices(token auth.Token, query url.Values) (devices []models.Device, err error, code int) {
	req, err := http.NewRequest("GET", this.config.DeviceRepoUrl+"/devices?"+query.Encode(), nil)
	if err != nil {
		debug.PrintStack()
		return devices, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return devices, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("error: status=%v; body=%v", resp.StatusCode, string(buf))
		return devices, err, resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		debug.PrintStack()
		return devices, err, http.StatusInternalServerError
	}
	return devices, nil, http.StatusOK
}

func (this *Com) GetDevice(token auth.Token, id string) (device models.Device, err error, code int) {
	err, code = getResourceFromService(token, this.config.DeviceRepoUrl+"/devices", id, &device)
	return
}

func (this *Com) GetDeviceByLocalId(token auth.Token, ownerId string, localid string) (device models.Device, err error, code int) {
	err, code = getResourceFromServiceWithQueryParam(token, this.config.DeviceRepoUrl+"/devices", localid, url.Values{"as": {"local_id"}, "owner_id": {ownerId}}, &device)
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
