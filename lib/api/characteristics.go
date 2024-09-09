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
	endpoints = append(endpoints, &CharacteristicsEndpoints{})
}

type CharacteristicsEndpoints struct{}

// Get godoc
// @Summary      get characteristic
// @Description  get characteristic
// @Tags         get, characteristics
// @Produce      json
// @Security Bearer
// @Param        id path string true "Characteristics Id"
// @Success      200 {object}  models.Characteristic
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /characteristics/{id} [GET]
func (this *CharacteristicsEndpoints) Get(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("GET /characteristics/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, errCode := control.ReadCharacteristic(token, id)
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
// @Summary      create characteristic
// @Description  create characteristic
// @Tags         create, characteristics
// @Produce      json
// @Security Bearer
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.Characteristic true "element"
// @Success      200 {object}  models.Characteristic
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /characteristics [POST]
func (this *CharacteristicsEndpoints) Create(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("POST /characteristics", func(writer http.ResponseWriter, request *http.Request) {
		characteristic := models.Characteristic{}
		err := json.NewDecoder(request.Body).Decode(&characteristic)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if characteristic.Id != "" {
			http.Error(writer, "body may not contain a preset id. please use the PUT method for updates", http.StatusBadRequest)
			return
		}

		options := model.CharacteristicUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishCharacteristicCreate(token, characteristic, options)
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
// @Summary      set characteristic
// @Description  set characteristic
// @Tags         set, characteristics
// @Produce      json
// @Security Bearer
// @Param        id path string true "Characteristic Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Param        message body models.Characteristic true "element"
// @Success      200 {object}  models.Characteristic
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /characteristics/{id} [PUT]
func (this *CharacteristicsEndpoints) Set(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("PUT /characteristics/{id}", func(writer http.ResponseWriter, request *http.Request) {
		characteristicId := request.PathValue("id")
		characteristic := models.Characteristic{}
		err := json.NewDecoder(request.Body).Decode(&characteristic)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.CharacteristicUpdateOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		result, err, errCode := control.PublishCharacteristicUpdate(token, characteristicId, characteristic, options)
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
// @Summary      delete characteristic
// @Description  delete characteristic
// @Tags         delete, characteristics
// @Produce      json
// @Security Bearer
// @Param        id path string true "Characteristic Id"
// @Param        wait query bool false "wait for done message in kafka before responding"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /characteristics/{id} [DELETE]
func (this *CharacteristicsEndpoints) Delete(config config.Config, router *http.ServeMux, control Controller) {
	router.HandleFunc("DELETE /characteristics/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		options := model.CharacteristicDeleteOptions{}
		if waitQueryParam := request.URL.Query().Get(WaitQueryParamName); waitQueryParam != "" {
			options.Wait, err = strconv.ParseBool(waitQueryParam)
			if err != nil {
				http.Error(writer, fmt.Sprintf("invalid %v query parameter %v", WaitQueryParamName, err.Error()), http.StatusBadRequest)
				return
			}
		}

		err, errCode := control.PublishCharacteristicDelete(token, id, options)
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
