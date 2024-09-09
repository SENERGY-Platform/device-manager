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
	"github.com/SENERGY-Platform/device-manager/lib/kafka/publisher"
	"github.com/SENERGY-Platform/models/go/models"
	"net/http"
	"net/http/httptest"
	"strconv"
)

type DeviceRepo struct {
	db              map[string]interface{}
	ts              *httptest.Server
	localIds        map[string]models.Device
	concepts        map[string]models.Concept
	characteristics map[string]models.Characteristic
	aspects         map[string]models.Aspect
	functions       map[string]models.Function
	deviceclasses   map[string]models.DeviceClass
	locations       map[string]models.Location
}

func NewDeviceRepo(producer interface {
	Subscribe(topic string, f func(msg []byte))
}) *DeviceRepo {
	repo := &DeviceRepo{
		db:              map[string]interface{}{},
		localIds:        map[string]models.Device{},
		concepts:        map[string]models.Concept{},
		characteristics: map[string]models.Characteristic{},
		aspects:         map[string]models.Aspect{},
		functions:       map[string]models.Function{},
		deviceclasses:   map[string]models.DeviceClass{},
		locations:       map[string]models.Location{},
	}
	producer.Subscribe(DtTopic, func(msg []byte) {
		cmd := publisher.DeviceTypeCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			for i, service := range cmd.DeviceType.Services {
				cmd.DeviceType.Services[i] = service
			}
			repo.db[cmd.Id] = cmd.DeviceType
		} else if cmd.Command == "DELETE" {
			delete(repo.db, cmd.Id)
		}
	})
	producer.Subscribe(DeviceTopic, func(msg []byte) {
		cmd := publisher.DeviceCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.localIds[cmd.Device.LocalId] = cmd.Device
			repo.db[cmd.Id] = cmd.Device
		} else if cmd.Command == "DELETE" {
			delete(repo.db, cmd.Id)
		}
	})

	producer.Subscribe(DeviceGroupTopic, func(msg []byte) {
		cmd := publisher.DeviceGroupCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.db[cmd.Id] = cmd.DeviceGroup
		} else if cmd.Command == "DELETE" {
			delete(repo.db, cmd.Id)
		}
	})

	producer.Subscribe(HubTopic, func(msg []byte) {
		cmd := publisher.HubCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.db[cmd.Id] = cmd.Hub
		} else if cmd.Command == "DELETE" {
			delete(repo.db, cmd.Id)
		}
	})
	producer.Subscribe(ProtocolTopic, func(msg []byte) {
		cmd := publisher.ProtocolCommand{}
		json.Unmarshal(msg, &cmd)
		if cmd.Command == "PUT" {
			repo.db[cmd.Id] = cmd.Protocol
		} else if cmd.Command == "DELETE" {
			delete(repo.db, cmd.Id)
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

	router := http.NewServeMux()

	router.HandleFunc("GET /device-groups/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		dt, ok := repo.db[id]
		if ok {
			json.NewEncoder(writer).Encode(dt)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /device-groups", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		group := models.DeviceGroup{}
		err = json.NewDecoder(request.Body).Decode(&group)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if group.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /device-types/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		dt, ok := repo.db[id]
		if ok {
			json.NewEncoder(writer).Encode(dt)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /device-types", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		dt := models.DeviceType{}
		err = json.NewDecoder(request.Body).Decode(&dt)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if dt.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		if len(dt.Services) == 0 {
			http.Error(writer, "expect at least one service", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		dt, ok := repo.db[id]
		if ok {
			json.NewEncoder(writer).Encode(dt)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /devices", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		device := models.Device{}
		err = json.NewDecoder(request.Body).Decode(&device)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if device.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		if device.LocalId == "" {
			http.Error(writer, "missing local id", http.StatusBadRequest)
			return
		}
		if device.DeviceTypeId == "" {
			http.Error(writer, "missing device-type id", http.StatusBadRequest)
			return
		}
		if existing, ok := repo.localIds[device.LocalId]; ok && existing.Id != device.Id {
			http.Error(writer, "expect local id to be globally unique", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /hubs/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		dt, ok := repo.db[id]
		if ok {
			json.NewEncoder(writer).Encode(dt)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /hubs", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		hub := models.Hub{}
		err = json.NewDecoder(request.Body).Decode(&hub)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if hub.Id == "" {
			http.Error(writer, "missing hub id", http.StatusBadRequest)
			return
		}
		if hub.Name == "" {
			http.Error(writer, "missing hub name", http.StatusBadRequest)
			return
		}

		if len(hub.DeviceLocalIds) == 1 {
			if _, ok := repo.localIds[hub.DeviceLocalIds[0]]; !ok {
				http.Error(writer, "unknown device local id", http.StatusBadRequest)
				return
			}
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("PUT /protocols", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		protocol := models.Protocol{}
		err = json.NewDecoder(request.Body).Decode(&protocol)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if protocol.Id == "" {
			http.Error(writer, "missing device id", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /aspects/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		aspect, ok := repo.aspects[id]
		if ok {
			json.NewEncoder(writer).Encode(aspect)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("DELETE /aspects/{id}", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("PUT /aspects", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		aspect := models.Aspect{}
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

	router.HandleFunc("GET /functions/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		function, ok := repo.functions[id]
		if ok {
			json.NewEncoder(writer).Encode(function)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /functions", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		function := models.Function{}
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

	router.HandleFunc("DELETE /functions/{id}", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /device-classes/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		deviceclass, ok := repo.deviceclasses[id]
		if ok {
			json.NewEncoder(writer).Encode(deviceclass)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /device-classes", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		deviceclass := models.DeviceClass{}
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

	router.HandleFunc("DELETE /device-classes/{id}", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("GET /locations/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		location, ok := repo.locations[id]
		if ok {
			json.NewEncoder(writer).Encode(location)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /locations", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		location := models.Location{}
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

	router.HandleFunc("GET /concepts/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		concept, ok := repo.concepts[id]
		if ok {
			json.NewEncoder(writer).Encode(concept)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("PUT /concepts", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		concept := models.Concept{}
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

	router.HandleFunc("DELETE /concepts/{id}", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("PUT /characteristics", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		characteristic := models.Characteristic{}
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

	router.HandleFunc("GET /characteristics/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		concept, ok := repo.characteristics[id]
		if ok {
			json.NewEncoder(writer).Encode(concept)
		} else {
			http.Error(writer, "404", 404)
		}
	})

	router.HandleFunc("DELETE /characteristics/{id}", func(writer http.ResponseWriter, request *http.Request) {
		dryRun, err := strconv.ParseBool(request.URL.Query().Get("dry-run"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if !dryRun {
			http.Error(writer, "only with query-parameter 'dry-run=true' allowed", http.StatusNotImplemented)
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	repo.ts = httptest.NewServer(router)

	return repo
}

func (this *DeviceRepo) Stop() {
	this.ts.Close()
}

func (this *DeviceRepo) Url() string {
	return this.ts.URL
}
