package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"runtime/debug"
)

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
