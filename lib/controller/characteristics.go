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
	"runtime/debug"
)

func (this *Controller) PublishCharacteristicCreate(token auth.Token, characteristic models.Characteristic, options model.CharacteristicUpdateOptions) (models.Characteristic, error, int) {
	if characteristic.Id != "" {
		return characteristic, errors.New("expect empty id"), http.StatusBadRequest
	}

	characteristic.GenerateId()
	err, code := this.com.ValidateCharacteristic(token, characteristic)
	if err != nil {
		return characteristic, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.CharacteristicTopic,
		ResourceId:   characteristic.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishCharacteristic(characteristic, token.GetUserId(), options.Wait)
	if err != nil {
		return characteristic, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return characteristic, err, http.StatusInternalServerError
	}

	return characteristic, nil, http.StatusOK
}

func (this *Controller) PublishCharacteristicUpdate(token auth.Token, characteristicId string, characteristic models.Characteristic, options model.CharacteristicUpdateOptions) (models.Characteristic, error, int) {
	if characteristic.Id != characteristicId {
		return characteristic, errors.New("characteristic id in body unequal to characteristic id in request endpoint"), http.StatusBadRequest
	}

	characteristic.GenerateId()

	if !token.IsAdmin() {
		err, code := this.com.PermissionCheckForCharacteristic(token, characteristicId, "w")
		if err != nil {
			debug.PrintStack()
			return characteristic, err, code
		}
	}
	err, code := this.com.ValidateCharacteristic(token, characteristic)
	if err != nil {
		debug.PrintStack()
		return characteristic, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.CharacteristicTopic,
		ResourceId:   characteristic.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishCharacteristic(characteristic, token.GetUserId(), options.Wait)
	if err != nil {
		debug.PrintStack()
		return characteristic, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return characteristic, err, http.StatusInternalServerError
	}

	return characteristic, nil, http.StatusOK
}

func (this *Controller) PublishCharacteristicDelete(token auth.Token, id string, options model.CharacteristicDeleteOptions) (error, int) {
	err, code := this.com.PermissionCheckForCharacteristic(token, id, "a")
	if err != nil {
		return err, code
	}
	err, code = this.com.ValidateCharacteristicDelete(token, id)
	if err != nil {
		return err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.CharacteristicTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishCharacteristicDelete(id, token.GetUserId(), options.Wait)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (this *Controller) ReadCharacteristic(token auth.Token, id string) (result models.Characteristic, err error, code int) {
	return this.com.GetCharacteristic(token, id)
}
