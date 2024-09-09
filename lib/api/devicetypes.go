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
	"strings"
)

func init() {
	endpoints = append(endpoints, &DeviceTypesEndpoints{})
}

type DeviceTypesEndpoints struct{}

// Get godoc
// @Summary      get device-type
// @Description  get device-type
// @Tags         get, device-types
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceType Id"
// @Success      200 {object}  models.DeviceType
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-types/{id} [GET]
func (this *DeviceTypesEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /device-types/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ReadDeviceType(token, id)
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
// @Summary      create device-type
// @Description  create device-type
// @Tags         create, device-types
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        distinct_attributes query string false "comma separated list of attribute keys; no other device-type with the same attribute key/value may exist"
// @Param        message body models.DeviceType true "element"
// @Success      200 {object}  models.DeviceType
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-types [POST]
func (this *DeviceTypesEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /device-types", func(writer http.ResponseWriter, request *http.Request) {
		devicetype := models.DeviceType{}
		err := json.NewDecoder(request.Body).Decode(&devicetype)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if devicetype.Id != "" {
			http.Error(writer, "body may not contain a preset id. please use the PUT method for updates", http.StatusBadRequest)
			return
		}

		distinctAttr := request.URL.Query().Get("distinct_attributes")
		if distinctAttr != "" {
			err = control.ValidateDistinctDeviceTypeAttributes(token, devicetype, strings.Split(distinctAttr, ","))
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		options := model.DeviceTypeUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceTypeCreate(token, devicetype, options)
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
// @Summary      set device-type
// @Description  set device-type
// @Tags         set, device-types
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceType Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        distinct_attributes query string false "comma separated list of attribute keys; no other device-type with the same attribute key/value may exist"
// @Param        message body models.DeviceType true "element"
// @Success      200 {object}  models.DeviceType
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-types/{id} [PUT]
func (this *DeviceTypesEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /device-types/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		devicetype := models.DeviceType{}
		err := json.NewDecoder(request.Body).Decode(&devicetype)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		distinctAttr := request.URL.Query().Get("distinct_attributes")
		if distinctAttr != "" {
			err = control.ValidateDistinctDeviceTypeAttributes(token, devicetype, strings.Split(distinctAttr, ","))
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		options := model.DeviceTypeUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceTypeUpdate(token, id, devicetype, options)
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
// @Summary      get device-type
// @Description  get device-type
// @Tags         get, device-types
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceType Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-types/{id} [DELETE]
func (this *DeviceTypesEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /device-types/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.DeviceTypeDeleteOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		err, errCode := control.PublishDeviceTypeDelete(token, id, options)
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
