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
	"log"
	"net/http"
	"strconv"
)

func init() {
	endpoints = append(endpoints, &DeviceGroupsEndpoints{})
}

type DeviceGroupsEndpoints struct{}

// Get godoc
// @Summary      get device-group
// @Description  get device-group
// @Tags         get, device-groups
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceGroup Id"
// @Success      200 {object}  models.DeviceGroup
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-groups/{id} [GET]
func (this *DeviceGroupsEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /device-groups/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ReadDeviceGroup(token, id)
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
// @Summary      create device-group
// @Description  create device-group
// @Tags         create, device-groups
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.DeviceGroup true "element"
// @Success      200 {object}  models.DeviceGroup
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-groups [POST]
func (this *DeviceGroupsEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /device-groups", func(writer http.ResponseWriter, request *http.Request) {
		deviceGroup := models.DeviceGroup{}
		err := json.NewDecoder(request.Body).Decode(&deviceGroup)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if deviceGroup.Id != "" {
			http.Error(writer, "device-group may not contain a preset id. please use PUT to update a device-group", http.StatusBadRequest)
			return
		}

		options := model.DeviceGroupUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceGroupCreate(token, deviceGroup, options)
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
// @Summary      set device-group
// @Description  set device-group
// @Tags         set, device-groups
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceGroup Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.DeviceGroup true "element"
// @Success      200 {object}  models.DeviceGroup
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-groups/{id} [PUT]
func (this *DeviceGroupsEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /device-groups/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		deviceGroup := models.DeviceGroup{}
		err := json.NewDecoder(request.Body).Decode(&deviceGroup)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.DeviceGroupUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceGroupUpdate(token, id, deviceGroup, options)
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
// @Summary      delete device-group
// @Description  delete device-group
// @Tags         delete, device-groups
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceGroup Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-groups/{id} [DELETE]
func (this *DeviceGroupsEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /device-groups/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.DeviceGroupDeleteOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		err, errCode := control.PublishDeviceGroupDelete(token, id, options)
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
