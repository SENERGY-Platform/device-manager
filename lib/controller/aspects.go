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

func (this *Controller) ReadAspect(token auth.Token, id string) (aspect models.Aspect, err error, code int) {
	return this.com.GetAspect(token, id)
}

func (this *Controller) PublishAspectCreate(token auth.Token, aspect models.Aspect, options model.AspectUpdateOptions) (models.Aspect, error, int) {
	if !token.IsAdmin() {
		return aspect, errors.New("access denied"), http.StatusForbidden
	}
	aspect.GenerateId()
	err, code := this.com.ValidateAspect(token, aspect)
	if err != nil {
		return aspect, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.AspectTopic,
		ResourceId:   aspect.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishAspect(aspect, token.GetUserId())
	if err != nil {
		return aspect, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return aspect, err, http.StatusInternalServerError
	}

	return aspect, nil, http.StatusOK
}

func (this *Controller) PublishAspectUpdate(token auth.Token, id string, aspect models.Aspect, options model.AspectUpdateOptions) (models.Aspect, error, int) {
	if !token.IsAdmin() {
		return aspect, errors.New("access denied"), http.StatusForbidden
	}
	if aspect.Id != id {
		return aspect, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	aspect.GenerateId()
	aspect.Id = id

	err, code := this.com.ValidateAspect(token, aspect)
	if err != nil {
		return aspect, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.AspectTopic,
		ResourceId:   aspect.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishAspect(aspect, token.GetUserId())
	if err != nil {
		return aspect, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return aspect, err, http.StatusInternalServerError
	}

	return aspect, nil, http.StatusOK
}

func (this *Controller) PublishAspectDelete(token auth.Token, id string, options model.AspectDeleteOptions) (error, int) {
	if !token.IsAdmin() {
		return errors.New("access denied"), http.StatusForbidden
	}
	err, code := this.com.ValidateAspectDelete(token, id)
	if err != nil {
		return err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.AspectTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishAspectDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
