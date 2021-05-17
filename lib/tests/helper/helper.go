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

package helper

import (
	"bytes"
	"encoding/json"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net"
	"net/http"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func Jwtdelete(jwt jwt_http_router.JwtImpersonate, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(jwt))
	resp, err = http.DefaultClient.Do(req)
	return
}

func JwtDeleteWithBody(jwt jwt_http_router.JwtImpersonate, url string, msg interface{}) (resp *http.Response, err error) {
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(msg)
	if err != nil {
		return resp, err
	}
	req, err := http.NewRequest("DELETE", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(jwt))
	resp, err = http.DefaultClient.Do(req)
	return
}

func Jwtget(jwt jwt_http_router.JwtImpersonate, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(jwt))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	return
}

func Jwtput(jwt jwt_http_router.JwtImpersonate, url string, msg interface{}) (resp *http.Response, err error) {
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(msg)
	if err != nil {
		return resp, err
	}
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(jwt))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	return
}

func Jwtpost(jwt jwt_http_router.JwtImpersonate, url string, msg interface{}) (resp *http.Response, err error) {
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(msg)
	if err != nil {
		return resp, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(jwt))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	return
}
