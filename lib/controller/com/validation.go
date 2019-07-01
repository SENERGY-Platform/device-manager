package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SmartEnergyPlatform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Com) ValidateDeviceType(jwt jwt_http_router.Jwt, dt model.DeviceType) (err error, code int) {
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(dt)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	for _, endpoint := range []string{
		this.config.SemanticDeviceRepoUrl + "/device-types",
		this.config.DeviceRepoUrl + "/device-types",
	} {
		req, err := http.NewRequest("HEAD", endpoint, b)
		if err != nil {
			return err, http.StatusInternalServerError
		}
		req.Header.Set("Authorization", string(jwt.Impersonate))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err, http.StatusInternalServerError
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			return errors.New(buf.String()), http.StatusBadRequest
		}
	}
	return nil, http.StatusOK
}
