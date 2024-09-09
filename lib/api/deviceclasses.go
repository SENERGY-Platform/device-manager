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
	endpoints = append(endpoints, &DeviceClassesEndpoints{})
}

type DeviceClassesEndpoints struct{}

// Get godoc
// @Summary      get device-class
// @Description  get device-class
// @Tags         get, device-classes
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceClass Id"
// @Success      200 {object}  models.DeviceClass
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-classes/{id} [GET]
func (this *DeviceClassesEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /device-classes/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ReadDeviceClass(token, id)
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
// @Summary      create device-class
// @Description  create device-class
// @Tags         create, device-classes
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.DeviceClass true "element"
// @Success      200 {object}  models.DeviceClass
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-classes [POST]
func (this *DeviceClassesEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /device-classes", func(writer http.ResponseWriter, request *http.Request) {
		deviceClass := models.DeviceClass{}
		err := json.NewDecoder(request.Body).Decode(&deviceClass)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if deviceClass.Id != "" {
			http.Error(writer, "body may not contain a preset id. please use the PUT method for updates", http.StatusBadRequest)
			return
		}

		options := model.DeviceClassUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceClassCreate(token, deviceClass, options)
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
// @Summary      set device-class
// @Description  set device-class
// @Tags         set, device-classes
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceClass Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.DeviceClass true "element"
// @Success      200 {object}  models.DeviceClass
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-classes/{id} [PUT]
func (this *DeviceClassesEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /device-classes/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		deviceClass := models.DeviceClass{}
		err := json.NewDecoder(request.Body).Decode(&deviceClass)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.DeviceClassUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishDeviceClassUpdate(token, id, deviceClass, options)
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
// @Summary      delete device-class
// @Description  delete device-class
// @Tags         delete, device-classes
// @Produce      json
// @Security Bearer
// @Param        id path string true "DeviceClass Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /device-classes/{id} [DELETE]
func (this *DeviceClassesEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /device-classes/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.DeviceClassDeleteOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		err, errCode := control.PublishDeviceClassDelete(token, id, options)
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
