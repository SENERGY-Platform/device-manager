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
	dt              map[string]model.DeviceType
	concepts        map[string]model.Concept
	characteristics map[string]model.Characteristic
	aspects         map[string]model.Aspect
	functions       map[string]model.Function
	deviceclasses   map[string]model.DeviceClass
	locations       map[string]model.Location
	ts              *httptest.Server
}

func NewSemanticRepo(producer interface {
	Subscribe(topic string, f func(msg []byte))
}) *SemanticRepo {
	repo := &SemanticRepo{
		dt:              map[string]model.DeviceType{},
		concepts:        map[string]model.Concept{},
		characteristics: map[string]model.Characteristic{},
		aspects:         map[string]model.Aspect{},
		functions:       map[string]model.Function{},
		deviceclasses:   map[string]model.DeviceClass{},
		locations:       map[string]model.Location{},
	}
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
	producer.Subscribe(ConceptTopic, func(msg []byte) {
		cmd := publisher.ConceptCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.concepts[cmd.Id] = cmd.Concept
		} else if cmd.Command == "DELETE" {
			delete(repo.concepts, cmd.Id)
		}
	})
	producer.Subscribe(CharacteristicTopic, func(msg []byte) {
		cmd := publisher.CharacteristicCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.characteristics[cmd.Id] = cmd.Characteristic
		} else if cmd.Command == "DELETE" {
			delete(repo.characteristics, cmd.Id)
		}
	})

	producer.Subscribe(AspectTopic, func(msg []byte) {
		cmd := publisher.AspectCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.aspects[cmd.Id] = cmd.Aspect
		} else if cmd.Command == "DELETE" {
			delete(repo.aspects, cmd.Id)
		}
	})

	producer.Subscribe(FunctionTopic, func(msg []byte) {
		cmd := publisher.FunctionCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.functions[cmd.Id] = cmd.Function
		} else if cmd.Command == "DELETE" {
			delete(repo.functions, cmd.Id)
		}
	})

	producer.Subscribe(DeviceClassTopic, func(msg []byte) {
		cmd := publisher.DeviceClassCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.deviceclasses[cmd.Id] = cmd.DeviceClass
		} else if cmd.Command == "DELETE" {
			delete(repo.deviceclasses, cmd.Id)
		}
	})

	producer.Subscribe(LocationTopic, func(msg []byte) {
		cmd := publisher.LocationCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.locations[cmd.Id] = cmd.Location
		} else if cmd.Command == "DELETE" {
			delete(repo.locations, cmd.Id)
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

	router.GET("/aspects/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		aspect, ok := repo.aspects[id]
		if ok {
			json.NewEncoder(writer).Encode(aspect)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.PUT("/aspects", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		aspect := model.Aspect{}
		err = json.NewDecoder(request.Body).Decode(&aspect)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if aspect.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.GET("/functions/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		function, ok := repo.functions[id]
		if ok {
			json.NewEncoder(writer).Encode(function)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.PUT("/functions", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		function := model.Function{}
		err = json.NewDecoder(request.Body).Decode(&function)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if function.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.GET("/device-classes/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		deviceclass, ok := repo.deviceclasses[id]
		if ok {
			json.NewEncoder(writer).Encode(deviceclass)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.PUT("/device-classes", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		deviceclass := model.DeviceClass{}
		err = json.NewDecoder(request.Body).Decode(&deviceclass)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if deviceclass.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.GET("/locations/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		location, ok := repo.locations[id]
		if ok {
			json.NewEncoder(writer).Encode(location)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.PUT("/locations", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		location := model.Location{}
		err = json.NewDecoder(request.Body).Decode(&location)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if location.Id == "" {
			http.Error(writer, "missing location id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.GET("/concepts/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		concept, ok := repo.concepts[id]
		if ok {
			json.NewEncoder(writer).Encode(concept)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.PUT("/concepts", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		concept := model.Concept{}
		err = json.NewDecoder(request.Body).Decode(&concept)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if concept.Id == "" {
			http.Error(writer, "missing concept id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.PUT("/characteristics", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		characteristic := model.Characteristic{}
		err = json.NewDecoder(request.Body).Decode(&characteristic)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if characteristic.Id == "" {
			http.Error(writer, "missing characteristic id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.GET("/characteristics/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		concept, ok := repo.characteristics[id]
		if ok {
			json.NewEncoder(writer).Encode(concept)
		} else {
			http.Error(writer, "404", 404)
		}
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
