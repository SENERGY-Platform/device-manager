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
	"log"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	endpoints = append(endpoints, &LocalDevicesEndpoints{})
}

type LocalDevicesEndpoints struct{}

// List godoc
// @Summary      list devices (local-id variant)
// @Description  list devices (local-id variant)
// @Tags         list, devices
// @Produce      json
// @Security Bearer
// @Param        ids query string false "comma separated list of local ids"
// @Param        owner_id query string false "defaults to requesting user; used in combination with local_id to find devices"
// @Param        limit query integer false "default 100, will be ignored if 'ids' is set"
// @Param        offset query integer false "default 0, will be ignored if 'ids' is set"
// @Param        search query string false "filter"
// @Param        sort query string false "default name.asc"
// @Param        device-type-ids query string false "filter; comma-seperated list"
// @Param        attr-keys query string false "filter; comma-seperated list; lists elements only if they have an attribute key that is in the given list"
// @Param        attr-values query string false "filter; comma-seperated list; lists elements only if they have an attribute value that is in the given list"
// @Param        connection-state query integer false "filter; valid values are 'online', 'offline' and an empty string for unknown states"
// @Success      200 {array}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /local-devices/{id} [GET]
func (this *LocalDevicesEndpoints) List(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /local-devices", func(writer http.ResponseWriter, request *http.Request) {
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
			result, err, errCode = control.ListDevicesByQuery(token, query)
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
}

// Get godoc
// @Summary      get device by local id
// @Description  get device by local id
// @Tags         get, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Local Id"
// @Param        owner_id query string false "defaults to requesting user; used in combination with id to find device"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /local-devices/{id} [GET]
func (this *LocalDevicesEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /local-devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
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
}

// Create godoc
// @Summary      create device (local-id variant)
// @Description  create device (local-id variant)
// @Tags         create, devices
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.Device true "element"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /local-devices/{id} [POST]
func (this *LocalDevicesEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /local-devices", func(writer http.ResponseWriter, request *http.Request) {
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
}

// Set godoc
// @Summary      set device (local-id variant)
// @Description  set device (local-id variant)
// @Tags         set, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Local Id"
// @Param        update-only-same-origin-attributes query string false "comma separated list; ensure that no attribute from another origin is overwritten"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.Device true "element"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /local-devices/{id} [PUT]
func (this *LocalDevicesEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /local-devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
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
}

// Delete godoc
// @Summary      delete device (local-id variant)
// @Description  delete device (local-id variant)
// @Tags         delete, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Local Id"
// @Param        owner_id query string false "defaults to requesting user; used in combination with id to find device"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /local-devices/{id} [DELETE]
func (this *LocalDevicesEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /local-devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
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
