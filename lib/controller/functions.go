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
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
	"net/http"
)

func (this *Controller) ReadFunction(token auth.Token, id string) (function models.Function, err error, code int) {
	return this.com.GetFunction(token, id)
}

func (this *Controller) PublishFunctionCreate(token auth.Token, function models.Function, options model.FunctionUpdateOptions) (models.Function, error, int) {
	if !token.IsAdmin() {
		return function, errors.New("access denied"), http.StatusForbidden
	}
	function.GenerateId()
	err, code := this.com.ValidateFunction(token, function)
	if err != nil {
		return function, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.FunctionTopic,
		ResourceId:   function.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishFunction(function, token.GetUserId(), options.Wait)
	if err != nil {
		return function, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return function, err, http.StatusInternalServerError
	}

	return function, nil, http.StatusOK
}

func (this *Controller) PublishFunctionUpdate(token auth.Token, id string, function models.Function, options model.FunctionUpdateOptions) (models.Function, error, int) {
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

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.FunctionTopic,
		ResourceId:   function.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishFunction(function, token.GetUserId(), options.Wait)
	if err != nil {
		return function, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return function, err, http.StatusInternalServerError
	}

	return function, nil, http.StatusOK
}

func (this *Controller) PublishFunctionDelete(token auth.Token, id string, options model.FunctionDeleteOptions) (error, int) {
	if !token.IsAdmin() {
		return errors.New("access denied"), http.StatusForbidden
	}
	err, code := this.com.ValidateFunctionDelete(token, id)
	if err != nil {
		return err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.FunctionTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishFunctionDelete(id, token.GetUserId(), options.Wait)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
