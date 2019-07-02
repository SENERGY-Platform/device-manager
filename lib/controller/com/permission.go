package com

import (
	"bytes"
	"encoding/json"
	"errors"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *Com) PermissionCheckForDeviceType(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) {
	return this.PermissionCheck(jwt, id, permission, this.config.DeviceTypeTopic)
}

func (this *Com) PermissionCheck(jwt jwt_http_router.Jwt, id string, permission string, resource string) (err error, code int) {
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/jwt/check/"+url.QueryEscape(resource)+"/"+url.QueryEscape(id)+"/"+permission+"/bool", nil)
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
		resp.Body.Close()
		log.Println("DEBUG: PermissionCheck()", buf.String())
		err = errors.New("access denied")
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}

	var ok bool
	err = json.NewDecoder(resp.Body).Decode(&ok)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	if !ok {
		return errors.New("access denied"), http.StatusForbidden
	}
	return
}
