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
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
	"net/http"
)

func (this *Controller) ReadProtocol(token auth.Token, id string) (protocol models.Protocol, err error, code int) {
	return this.com.GetProtocol(token, id)
}

func (this *Controller) PublishProtocolCreate(token auth.Token, protocol models.Protocol, options model.ProtocolUpdateOptions) (models.Protocol, error, int) {
	if !token.IsAdmin() {
		return protocol, errors.New("access denied"), http.StatusForbidden
	}
	protocol.GenerateId()
	err, code := this.com.ValidateProtocol(token, protocol)
	if err != nil {
		return protocol, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.ProtocolTopic,
		ResourceId:   protocol.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishProtocol(protocol, token.GetUserId())
	if err != nil {
		return protocol, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return protocol, err, http.StatusInternalServerError
	}

	return protocol, nil, http.StatusOK
}

func (this *Controller) PublishProtocolUpdate(token auth.Token, id string, protocol models.Protocol, options model.ProtocolUpdateOptions) (models.Protocol, error, int) {
	if !token.IsAdmin() {
		return protocol, errors.New("access denied"), http.StatusForbidden
	}
	if protocol.Id != id {
		return protocol, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	//replace sub ids and create new ones for new sub elements
	protocol.GenerateId()
	protocol.Id = id

	err, code := this.com.ValidateProtocol(token, protocol)
	if err != nil {
		return protocol, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.ProtocolTopic,
		ResourceId:   protocol.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishProtocol(protocol, token.GetUserId())
	if err != nil {
		return protocol, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return protocol, err, http.StatusInternalServerError
	}

	return protocol, nil, http.StatusOK
}

func (this *Controller) PublishProtocolDelete(token auth.Token, id string, options model.ProtocolDeleteOptions) (error, int) {
	if !token.IsAdmin() {
		return errors.New("access denied"), http.StatusForbidden
	}
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.ProtocolTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err := this.publisher.PublishProtocolDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}
