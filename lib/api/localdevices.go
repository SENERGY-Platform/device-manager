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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/jwt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	endpoints = append(endpoints, LocalDevicesEndpoints)
}

func LocalDevicesEndpoints(config config.Config, control Controller, router *httprouter.Router) {
	resource := "/local-devices"

	//query-parameter:
	//		- limit: number; default 100, will be ignored if 'ids' is set
	//		- offset: number; default 0, will be ignored if 'ids' is set
	//		- ids: filter by comma seperated id list
	router.GET(resource, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		token, err := jwt.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		ownerId := request.URL.Query().Get("owner_id")
		if ownerId == "" {
			ownerId = token.GetUserId()
		}

		result := []models.Device{}

		query := request.URL.Query()
		idsStr := query.Get("ids")
		if idsStr != "" {
			localIds := strings.Split(idsStr, ",")
			for _, localId := range localIds {
				device, err, errCode := control.ReadDeviceByLocalId(token, ownerId, localId)
				if err == nil {
					result = append(result, device)
				}
				if err != nil && errCode >= 500 {
					http.Error(writer, err.Error(), errCode)
					return
				}
			}
		} else {
			var errCode int
			result, err, errCode = control.ListDevices(token, query)
			if err != nil {
				http.Error(writer, err.Error(), errCode)
				return
			}
		}

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})

	router.GET(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		ownerId := request.URL.Query().Get("owner_id")
		if ownerId == "" {
			ownerId = token.GetUserId()
		}
		result, err, errCode := control.ReadDeviceByLocalId(token, ownerId, id)
		if err != nil {
			http.Error(writer, err.Error(), errCode)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})

	router.POST(resource, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		device := models.Device{}
		err := json.NewDecoder(request.Body).Decode(&device)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if device.Id != "" {
			http.Error(writer, "body may not contain a preset id. please use the PUT method for updates", http.StatusBadRequest)
			return
		}

		options := model.DeviceCreateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceCreate(token, device, options)
		if err != nil {
			http.Error(writer, err.Error(), errCode)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})

	router.PUT(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		ownerId := token.GetUserId()
		id, err, errCode := control.DeviceLocalIdToId(token, ownerId, id)
		if err != nil {
			http.Error(writer, err.Error(), errCode)
			return
		}

		device := models.Device{}
		err = json.NewDecoder(request.Body).Decode(&device)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if device.Id != "" && device.Id != id {
			http.Error(writer, "device contains a different id then the id from the url", http.StatusBadRequest)
			return
		}
		device.Id = id

		options := model.DeviceUpdateOptions{}
		if request.URL.Query().Has(UpdateOnlySameOriginAttributesKey) {
			temp := request.URL.Query().Get(UpdateOnlySameOriginAttributesKey)
			options.UpdateOnlySameOriginAttributes = strings.Split(temp, ",")
		}

		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceUpdate(token, id, device, options)
		if err != nil {
			http.Error(writer, err.Error(), errCode)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})

	router.DELETE(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		ownerId := request.URL.Query().Get("owner_id")
		if ownerId == "" {
			ownerId = token.GetUserId()
		}
		id, err, errCode := control.DeviceLocalIdToId(token, ownerId, id)
		if err != nil {
			http.Error(writer, err.Error(), errCode)
			return
		}

		options := model.DeviceDeleteOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		err, errCode = control.PublishDeviceDelete(token, id, options)
		if err != nil {
			http.Error(writer, err.Error(), errCode)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(true)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})
}
