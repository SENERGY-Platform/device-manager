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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *Com) GetTechnicalDeviceType(jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int) {
	return this.getDeviceFromService(this.config.DeviceRepoUrl, jwt, id)
}

func (this *Com) GetSemanticDeviceType(jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int) {
	return this.getDeviceFromService(this.config.SemanticRepoUrl, jwt, id)
}

func (this *Com) getDeviceFromService(service string, jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int) {
	req, err := http.NewRequest("GET", service+"/device-types/"+url.PathEscape(id), nil)
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", string(jwt.Impersonate))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return dt, errors.New(buf.String()), resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&dt)
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Com) ValidateDeviceType(jwt jwt_http_router.Jwt, dt model.DeviceType) (err error, code int) {
	for _, endpoint := range []string{
		this.config.SemanticRepoUrl + "/device-types?dry-run=true",
		this.config.DeviceRepoUrl + "/device-types?dry-run=true",
	} {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(dt)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		req, err := http.NewRequest("PUT", endpoint, b)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		req.Header.Set("Authorization", string(jwt.Impersonate))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			return errors.New(buf.String()), resp.StatusCode
		}
	}
	return nil, http.StatusOK
}
