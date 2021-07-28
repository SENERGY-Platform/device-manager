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

package controller

import (
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"net/http"
)

func (this *Controller) ReadFunction(token auth.Token, id string) (function model.Function, err error, code int) {
	return this.com.GetFunction(token, id)
}

func (this *Controller) PublishFunctionCreate(token auth.Token, function model.Function) (model.Function, error, int) {
	if !token.IsAdmin() {
		return function, errors.New("access denied"), http.StatusForbidden
	}
	function.GenerateId()
	err, code := this.com.ValidateFunction(token, function)
	if err != nil {
		return function, err, code
	}
	err = this.publisher.PublishFunction(function, token.GetUserId())
	if err != nil {
		return function, err, http.StatusInternalServerError
	}
	return function, nil, http.StatusOK
}

func (this *Controller) PublishFunctionUpdate(token auth.Token, id string, function model.Function) (model.Function, error, int) {
	if !token.IsAdmin() {
		return function, errors.New("access denied"), http.StatusForbidden
	}
	if function.Id != id {
		return function, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	function.GenerateId()
	function.Id = id

	err, code := this.com.ValidateFunction(token, function)
	if err != nil {
		return function, err, code
	}
	err = this.publisher.PublishFunction(function, token.GetUserId())
	if err != nil {
		return function, err, http.StatusInternalServerError
	}
	return function, nil, http.StatusOK
}

func (this *Controller) PublishFunctionDelete(token auth.Token, id string) (error, int) {
	if !token.IsAdmin() {
		return errors.New("access denied"), http.StatusForbidden
	}
	err := this.publisher.PublishFunctionDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
