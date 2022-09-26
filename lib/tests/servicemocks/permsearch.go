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

package servicemocks

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
)

type PermSearch struct {
	ts *httptest.Server
}

func NewPermSearch() *PermSearch {
	repo := &PermSearch{}

	router := httprouter.New()

	router.GET("/jwt/check/:resource/:id/:permission/bool", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		json.NewEncoder(writer).Encode(true)
	})

	router.GET("/v3/resources/:resource/:id/access", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		json.NewEncoder(writer).Encode(true)
	})

	router.GET("/jwt/select/devices/device_type_id/:id/x", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		json.NewEncoder(writer).Encode([]interface{}{})
	})

	router.POST("/v3/query", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		message := com.QueryMessage{}
		err := json.NewDecoder(request.Body).Decode(&message)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}
		if message.CheckIds == nil {
			http.Error(writer, "not implemented", 500)
			return
		}
		resp := map[string]bool{}
		for _, id := range message.CheckIds.Ids {
			resp[id] = true
		}
		json.NewEncoder(writer).Encode(resp)
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
