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
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/http/httptest"
	"strconv"
)

type SemanticRepo struct {
	dt map[string]model.DeviceType
	ts *httptest.Server
}

func NewSemanticRepo(producer interface {
	Subscribe(topic string, f func(msg []byte))
}) *SemanticRepo {
	repo := &SemanticRepo{dt: map[string]model.DeviceType{}}
	producer.Subscribe(DtTopic, func(msg []byte) {
		cmd := publisher.DeviceTypeCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			for i, service := range cmd.DeviceType.Services {
				service.ProtocolId = ""
				cmd.DeviceType.Services[i] = service
			}
			repo.dt[cmd.Id] = cmd.DeviceType
		} else if cmd.Command == "DELETE" {
			delete(repo.dt, cmd.Id)
		}
	})

	router := jwt_http_router.New(jwt_http_router.JwtConfig{ForceAuth: true, ForceUser: true})

	router.GET("/device-types/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		dt, ok := repo.dt[id]
		if ok {
			json.NewEncoder(writer).Encode(dt)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.PUT("/device-types", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		dt := model.DeviceType{}
		err = json.NewDecoder(request.Body).Decode(&dt)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if dt.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	repo.ts = httptest.NewServer(router)

	return repo
}

func (this *SemanticRepo) Stop() {
	this.ts.Close()
}

func (this *SemanticRepo) Url() string {
	return this.ts.URL
}
