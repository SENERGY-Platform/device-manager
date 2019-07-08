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

package mock

import (
	"encoding/json"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/http/httptest"
)

type PermSearch struct {
	ts *httptest.Server
}

func NewPermSearch() *PermSearch {
	repo := &PermSearch{}

	router := jwt_http_router.New(jwt_http_router.JwtConfig{ForceAuth: true, ForceUser: true})

	router.GET("/jwt/check/:resource/:id/:permission/bool", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		json.NewEncoder(writer).Encode(true)
	})

	repo.ts = httptest.NewServer(router)

	return repo
}

func (this *PermSearch) Stop() {
	this.ts.Close()
}

func (this *PermSearch) Url() string {
	return this.ts.URL
}
