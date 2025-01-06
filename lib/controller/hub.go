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
	"fmt"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
	"log"
	"net/http"
	"runtime/debug"
	"slices"
)

func (this *Controller) ReadHub(token auth.Token, id string) (hub models.Hub, err error, code int) {
	return this.com.GetHub(token, id)
}

func (this *Controller) PublishHubCreate(token auth.Token, hubEdit models.Hub, options model.HubUpdateOptions) (models.Hub, error, int) {
	hub, err, code := this.completeHub(token, hubEdit)
	if err != nil {
		return hub, err, code
	}
	hub.GenerateId()
	if hub.OwnerId != "" && hub.OwnerId != token.GetUserId() {
		return hub, errors.New("new devices must be initialised with the requesting user as owner-id"), http.StatusBadRequest
	}
	hub.OwnerId = token.GetUserId()
	err, code = this.com.ValidateHub(token, hub)
	if err != nil {
		return hub, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.HubTopic,
		ResourceId:   hub.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishHub(hub, token.GetUserId())
	if err != nil {
		return hub, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return hub, err, http.StatusInternalServerError
	}

	return hub, nil, http.StatusOK
}

func (this *Controller) PublishHubUpdate(token auth.Token, id string, hubEdit models.Hub, options model.HubUpdateOptions) (models.Hub, error, int) {
	if hubEdit.Id != id {
		return models.Hub{}, errors.New("hub id in body unequal to hub id in request endpoint"), http.StatusBadRequest
	}

	hub, err, code := this.completeHub(token, hubEdit)
	if err != nil {
		return hub, err, code
	}

	if !token.IsAdmin() {
		err, code := this.com.PermissionCheckForHub(token, id, "w")
		if err != nil {
			return hub, err, code
		}
	}

	var original models.Hub
	var exists bool
	original, err, code = this.com.GetHub(token, hub.Id)
	if err != nil && code != http.StatusNotFound {
		return hub, err, code
	}
	if err != nil {
		err, code = nil, 200
		exists = false
	} else {
		exists = true
	}

	//set device owner-id if none is given
	//prefer existing owner, fallback to requesting user
	if hub.OwnerId == "" {
		hub.OwnerId = original.OwnerId //may be empty for new devices
	}
	if hub.OwnerId == "" {
		hub.OwnerId = token.GetUserId()
	}

	if exists && original.OwnerId != hub.OwnerId && !token.IsAdmin() {
		err, code := this.com.PermissionCheckForHub(token, hub.Id, "a")
		if err != nil {
			if code == http.StatusForbidden {
				return hub, fmt.Errorf("only admins may set new owner: %w", err), http.StatusBadRequest
			} else {
				return hub, err, code
			}
		}
	}

	err, code = this.com.ValidateHub(token, hub)
	if err != nil {
		return hub, err, code
	}

	rights, err, code := this.com.GetResourceRights(token, this.config.HubTopic, hub.Id)
	if err != nil && code != http.StatusNotFound {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return hub, err, code
	}

	//new device owner-id must be existing admin user (ignore for new devices or devices with unchanged owner)
	if code != http.StatusNotFound && hub.OwnerId != original.OwnerId && !rights.UserPermissions[hub.OwnerId].Administrate {
		return hub, errors.New("new owner must have existing user admin rights"), http.StatusBadRequest
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.HubTopic,
		ResourceId:   hub.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishHub(hub, hub.OwnerId)
	if err != nil {
		return hub, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return hub, err, http.StatusInternalServerError
	}

	return hub, nil, http.StatusOK
}

func (this *Controller) PublishHubDelete(token auth.Token, id string, options model.HubDeleteOptions) (error, int) {
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	err, code := this.com.PermissionCheckForHub(token, id, "a")
	if err != nil {
		return err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.HubTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishHubDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

type IdWrapper struct {
	Id string `json:"id"`
}

func (this *Controller) completeHub(token auth.Token, edit models.Hub) (result models.Hub, err error, code int) {
	result = edit
	if result.DeviceLocalIds == nil {
		result.DeviceLocalIds = []string{}
	}
	if result.DeviceIds == nil {
		result.DeviceIds = []string{}
	}

	if len(edit.DeviceLocalIds) > 0 {
		devices, err, code := this.com.ListDevices(token.Jwt(), client.DeviceListOptions{LocalIds: edit.DeviceLocalIds})
		if err != nil {
			return result, err, code
		}
		if len(devices) != len(edit.DeviceLocalIds) {
			return result, errors.New("not all local device ids found"), http.StatusBadRequest
		}
		for _, device := range devices {
			if !slices.Contains(result.DeviceLocalIds, device.LocalId) {
				result.DeviceLocalIds = append(result.DeviceLocalIds, device.LocalId)
			}
			if !slices.Contains(result.DeviceIds, device.Id) {
				result.DeviceIds = append(result.DeviceIds, device.Id)
			}
		}
	}

	if len(edit.DeviceIds) > 0 {
		devices, err, code := this.com.ListDevices(token.Jwt(), client.DeviceListOptions{Ids: edit.DeviceIds})
		if err != nil {
			return result, err, code
		}
		if len(devices) != len(edit.DeviceIds) {
			return result, errors.New("not all device ids found"), http.StatusBadRequest
		}
		for _, device := range devices {
			if !slices.Contains(result.DeviceLocalIds, device.LocalId) {
				result.DeviceLocalIds = append(result.DeviceLocalIds, device.LocalId)
			}
			if !slices.Contains(result.DeviceIds, device.Id) {
				result.DeviceIds = append(result.DeviceIds, device.Id)
			}
		}
	}

	if len(result.DeviceLocalIds) != len(result.DeviceIds) {
		return result, errors.New("DeviceLocalIds DeviceIds count mismatch"), http.StatusBadRequest
	}
	return result, err, code
}
