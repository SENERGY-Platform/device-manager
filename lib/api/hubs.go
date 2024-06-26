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
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func init() {
	endpoints = append(endpoints, HubsEndpoints)
}

func HubsEndpoints(config config.Config, control Controller, router *httprouter.Router) {
	resource := "/hubs"

	router.GET(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
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

	router.HEAD(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		_, _, errCode := control.ReadHub(token, id)
		writer.WriteHeader(errCode)
		return
	})

	router.POST(resource, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	router.PUT(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
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

	router.PUT(resource+"/:id/name", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
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

	router.DELETE(resource+"/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
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
