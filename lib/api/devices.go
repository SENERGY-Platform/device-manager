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
	endpoints = append(endpoints, &DevicesEndpoints{})
}

type DevicesEndpoints struct{}

const UpdateOnlySameOriginAttributesKey = "update-only-same-origin-attributes"
const DisplayNameAttributeKey = "shared/nickname"
const DisplayNameAttributeOrigin = "shared"

const WaitQueryParamName = "wait"

// List godoc
// @Summary      list devices
// @Description  list devices
// @Tags         list, devices
// @Produce      json
// @Security Bearer
// @Param        limit query integer false "default 100, will be ignored if 'ids' is set"
// @Param        offset query integer false "default 0, will be ignored if 'ids' is set"
// @Param        search query string false "filter"
// @Param        sort query string false "default name.asc"
// @Param        ids query string false "filter; ignores limit/offset; comma-seperated list"
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
// @Router       /devices [GET]
func (this *DevicesEndpoints) List(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /devices", func(writer http.ResponseWriter, request *http.Request) {
		token, err := jwt.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ListDevicesByQuery(token, request.URL.Query())
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

// Get godoc
// @Summary      get device
// @Description  get device
// @Tags         get, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Id"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /devices/{id} [GET]
func (this *DevicesEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ReadDevice(token, id)
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
// @Summary      create device
// @Description  create device
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
// @Router       /devices [POST]
func (this *DevicesEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /devices", func(writer http.ResponseWriter, request *http.Request) {
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
// @Summary      set device
// @Description  set device; admins may create new devices but only without using the UpdateOnlySameOriginAttributesKey query parameter
// @Tags         set, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        update-only-same-origin-attributes query string false "comma separated list; ensure that no attribute from another origin is overwritten"
// @Param        message body models.Device true "element"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /devices/{id} [PUT]
func (this *DevicesEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
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

// SetAttributes godoc
// @Summary      set device attributes
// @Description  set device attributes
// @Tags         set, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        update-only-same-origin-attributes query string false "comma separated list; ensure that no attribute from another origin is overwritten"
// @Param        message body []models.Attribute true "attributes"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /devices/{id}/attributes [PUT]
func (this *DevicesEndpoints) SetAttributes(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /devices/{id}/attributes", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		attributes := []models.Attribute{}
		err := json.NewDecoder(request.Body).Decode(&attributes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

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

		device, err, errCode := control.ReadDevice(token, id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		device.Attributes = attributes

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

// SetDisplayName godoc
// @Summary      set device display name
// @Description  set device display name
// @Tags         set, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body string true "display name"
// @Success      200 {object}  models.Device
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /devices/{id}/display_name [PUT]
func (this *DevicesEndpoints) SetDisplayName(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /devices/{id}/display_name", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		displayName := ""

		err := json.NewDecoder(request.Body).Decode(&displayName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		device, err, errCode := control.ReadDevice(token, id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		displayNameAttrFound := false
		for i, attr := range device.Attributes {
			if attr.Key == DisplayNameAttributeKey {
				attr.Value = displayName
				device.Attributes[i] = attr
				displayNameAttrFound = true
			}
		}
		if !displayNameAttrFound {
			device.Attributes = append(device.Attributes, models.Attribute{Key: DisplayNameAttributeKey, Value: displayName, Origin: DisplayNameAttributeOrigin})
		}

		options := model.DeviceUpdateOptions{}

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
// @Summary      delete device
// @Description  delete device
// @Tags         delete, devices
// @Produce      json
// @Security Bearer
// @Param        id path string true "Device Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /devices/{id} [DELETE]
func (this *DevicesEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /devices/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
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

		err, errCode := control.PublishDeviceDelete(token, id, options)
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

// DeleteMany godoc
// @Summary      delete multiple devices
// @Description  delete multiple devices
// @Tags         delete, devices
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body []string true "ids to be deleted"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /devices [DELETE]
func (this *DevicesEndpoints) DeleteMany(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /devices", func(writer http.ResponseWriter, request *http.Request) {
		ids := []string{}
		err := json.NewDecoder(request.Body).Decode(&ids)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
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

		for i, id := range ids {
			if i < len(ids)-1 {
				err, errCode := control.PublishDeviceDelete(token, id, model.DeviceDeleteOptions{})
				if err != nil {
					http.Error(writer, err.Error(), errCode)
					return
				}
			} else {
				err, errCode := control.PublishDeviceDelete(token, id, options)
				if err != nil {
					http.Error(writer, err.Error(), errCode)
					return
				}
			}
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(true)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
		return
	})
}
