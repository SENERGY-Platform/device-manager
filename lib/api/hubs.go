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
	endpoints = append(endpoints, &HubsEndpoints{})
}

type HubsEndpoints struct{}

// Get godoc
// @Summary      get hub
// @Description  get hub
// @Tags         get, hubs
// @Produce      json
// @Security Bearer
// @Param        id path string true "Hub Id"
// @Success      200 {object}  models.Hub
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /hubs/{id} [GET]
func (this *HubsEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /hubs/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ReadHub(token, id)
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

// Head godoc
// @Summary      head hub
// @Description  head hub
// @Tags         head, hubs
// @Produce      json
// @Security Bearer
// @Param        id path string true "Hub Id"
// @Param        message body models.Hub true "element"
// @Success      200 {object}  models.Hub
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /hubs/{id} [HEAD]
func (this *HubsEndpoints) Head(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("HEAD /hubs/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		_, _, errCode := control.ReadHub(token, id)
		writer.WriteHeader(errCode)
		return
	})
}

// Create godoc
// @Summary      create hub
// @Description  create hub
// @Tags         create, hubs
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.Hub true "element"
// @Success      200 {object}  models.Hub
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /hubs [POST]
func (this *HubsEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /hubs", func(writer http.ResponseWriter, request *http.Request) {
		hub := models.Hub{}
		err := json.NewDecoder(request.Body).Decode(&hub)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if hub.Id != "" {
			http.Error(writer, "body may not contain a preset id. please use the PUT method for updates", http.StatusBadRequest)
			return
		}

		options := model.HubUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishHubCreate(token, hub, options)
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
// @Summary      set hub
// @Description  set hub
// @Tags         set, hubs
// @Produce      json
// @Security Bearer
// @Param        id path string true "Hub Id"
// @Param        user_id query string false "only admins may set user_id; overwrites hub.OwnerId; defaults to existing hub.OwnerId and falls back to user-id of requesting user if hub does not exist"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.Hub true "element"
// @Success      200 {object}  models.Hub
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /hubs/{id} [PUT]
func (this *HubsEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /hubs/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		userId := request.URL.Query().Get("user_id")
		hub := models.Hub{}
		err := json.NewDecoder(request.Body).Decode(&hub)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if userId != "" && !token.IsAdmin() {
			http.Error(writer, "only admins may set user_id", http.StatusForbidden)
			return
		}

		options := model.HubUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishHubUpdate(token, id, userId, hub, options)
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

// SetName godoc
// @Summary      set hub name
// @Description  set hub name
// @Tags         set, hubs
// @Produce      json
// @Security Bearer
// @Param        id path string true "Hub Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body string true "name"
// @Success      200 {object}  models.Hub
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /hubs/{id}/name [PUT]
func (this *HubsEndpoints) SetName(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /hubs/{id}/name", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		name := ""
		err := json.NewDecoder(request.Body).Decode(&name)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		hub, err, code := control.ReadHub(token, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		hub.Name = name

		options := model.HubUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishHubUpdate(token, id, token.GetUserId(), hub, options)
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
// @Summary      delete hub
// @Description  delete hub
// @Tags         delete, hubs
// @Produce      json
// @Security Bearer
// @Param        id path string true "Hub Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /hubs/{id} [DELETE]
func (this *HubsEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /hubs/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.HubDeleteOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		err, errCode := control.PublishHubDelete(token, id, options)
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
