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
	"log"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
)

type PermSearch struct {
	ts *httptest.Server
}

func NewPermSearch() *PermSearch {
	repo := &PermSearch{}

	router := http.NewServeMux()

	router.HandleFunc("GET /jwt/check/{resource}/{id}/:permission/bool", func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(true)
	})

	router.HandleFunc("GET /v3/resources/{resource}/{id}/access", func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(true)
	})

	router.HandleFunc("GET /v3/resources/{resource}", func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode([]interface{}{})
	})

	router.HandleFunc("POST /v3/query", func(writer http.ResponseWriter, request *http.Request) {
		message := com.QueryMessage{}
		err := json.NewDecoder(request.Body).Decode(&message)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}
		if message.CheckIds != nil {
			resp := map[string]bool{}
			for _, id := range message.CheckIds.Ids {
				resp[id] = true
			}
			json.NewEncoder(writer).Encode(resp)
			return
		} else if message.Find != nil {
			json.NewEncoder(writer).Encode([]map[string]interface{}{})
			return
		} else if message.ListIds != nil {
			if len(message.ListIds.Ids) == 1 {
				//used for Com.GetResourceOwner()
				json.NewEncoder(writer).Encode([]map[string]interface{}{{"creator": "testowner"}})
			} else {
				json.NewEncoder(writer).Encode([]map[string]interface{}{})
			}
			return
		} else {
			temp, _ := json.Marshal(message)
			log.Println("ERROR: ", string(temp))
			debug.PrintStack()
			http.Error(writer, "not implemented", 500)
			return
		}
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
