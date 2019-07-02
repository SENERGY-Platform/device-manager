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
