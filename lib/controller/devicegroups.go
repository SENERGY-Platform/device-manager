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
	"log"
	"net/http"
	"runtime/debug"
)

func (this *Controller) ReadDeviceGroup(token auth.Token, id string) (dt models.DeviceGroup, err error, code int) {
	return this.com.GetTechnicalDeviceGroup(token, id)
}

func (this *Controller) PublishDeviceGroupCreate(token auth.Token, dg models.DeviceGroup, options model.DeviceGroupUpdateOptions) (result models.DeviceGroup, err error, code int) {
	if dg.Id != "" {
		return result, errors.New("expect empty device-group id"), http.StatusBadRequest
	}

	dg.GenerateId()
	dg.SetShortCriteria()
	dg.DeviceIds, err = this.filterInvalidDeviceIds(token, dg.DeviceIds, "r")
	if err != nil {
		return dg, err, http.StatusInternalServerError
	}
	err, code = this.com.ValidateDeviceGroup(token, dg)
	if err != nil {
		return dg, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceGroupTopic,
		ResourceId:   dg.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishDeviceGroup(dg, token.GetUserId())
	if err != nil {
		return dg, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return dg, err, http.StatusInternalServerError
	}

	return dg, nil, http.StatusOK
}

func (this *Controller) PublishDeviceGroupUpdate(token auth.Token, id string, dg models.DeviceGroup, options model.DeviceGroupUpdateOptions) (result models.DeviceGroup, err error, code int) {
	if dg.Id != id {
		return dg, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	dg.GenerateId()
	dg.SetShortCriteria()

	dg.DeviceIds, err = this.filterInvalidDeviceIds(token, dg.DeviceIds, "r")
	if err != nil {
		return dg, err, http.StatusInternalServerError
	}

	if !token.IsAdmin() {
		err, code := this.com.PermissionCheckForDeviceGroup(token, id, "w")
		if err != nil {
			debug.PrintStack()
			return dg, err, code
		}
	}
	err, code = this.com.ValidateDeviceGroup(token, dg)
	if err != nil {
		debug.PrintStack()
		return dg, err, code
	}

	//ensure retention of original owner
	owner, found, err := this.com.GetResourceOwner(token, this.config.DeviceGroupTopic, dg.Id, "w")
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return dg, err, http.StatusInternalServerError
	}
	if !found || owner == "" {
		owner = token.GetUserId()
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceGroupTopic,
		ResourceId:   dg.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishDeviceGroup(dg, owner)
	if err != nil {
		debug.PrintStack()
		return dg, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return dg, err, http.StatusInternalServerError
	}

	return dg, nil, http.StatusOK
}

func (this *Controller) PublishDeviceGroupDelete(token auth.Token, id string, options model.DeviceGroupDeleteOptions) (error, int) {
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	err, code := this.com.PermissionCheckForDeviceGroup(token, id, "a")
	if err != nil {
		return err, code
	}

	err, code = this.com.ValidateDeviceGroupDelete(token, id)
	if err != nil {
		return err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceGroupTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishDeviceGroupDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (this *Controller) filterInvalidDeviceIds(token auth.Token, ids []string, rights string) (result []string, err error) {
	deviceIsAccessible, err, _ := this.com.PermissionCheckForDeviceList(token, ids, rights)
	if err != nil {
		return result, err
	}
	result = []string{}
	for _, id := range ids {
		if deviceIsAccessible[id] {
			result = append(result, id)
		} else {
			log.Println("WARNING: remove device from device-group because its inaccessible", id)
		}
	}
	return result, nil
}
