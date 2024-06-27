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
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
	"strings"
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

func (this *Controller) PublishHubUpdate(token auth.Token, id string, userId string, hubEdit models.Hub, options model.HubUpdateOptions) (models.Hub, error, int) {
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

	rights, found, err := this.com.GetResourceRights(token, this.config.DeviceTopic, hub.Id, "w")
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return hub, err, http.StatusInternalServerError
	}

	//new device owner-id must be existing admin user (ignore for new devices or devices with unchanged owner)
	if found && hub.OwnerId != original.OwnerId && !slices.Contains(rights.PermissionHolders.AdminUsers, hub.OwnerId) {
		return hub, errors.New("new owner must have existing user admin rights"), http.StatusBadRequest
	}

	//ensure retention of original creator
	creator, found, err := this.com.GetResourceOwner(token, this.config.HubTopic, hub.Id, "w")
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return hub, err, http.StatusInternalServerError
	}
	if found && creator != "" {
		userId = creator
	}
	if userId == "" {
		userId = token.GetUserId()
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.HubTopic,
		ResourceId:   hub.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishHub(hub, userId)
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
	idWrapperList := []IdWrapper{}
	err, code = this.com.QueryPermissionsSearch(token.Jwt(), com.QueryMessage{
		Resource: "devices",
		Find: &com.QueryFind{
			QueryListCommons: com.QueryListCommons{
				Limit:  len(edit.DeviceLocalIds),
				Offset: 0,
				Rights: "r",
			},
			Filter: &com.Selection{
				Condition: com.ConditionConfig{
					Feature:   "features.local_id",
					Operation: com.QueryAnyValueInFeatureOperation,
					Value:     edit.DeviceLocalIds,
				},
			},
		},
	}, &idWrapperList)
	if err != nil {
		return result, err, code
	}
	for _, id := range idWrapperList {
		if !slices.Contains(result.DeviceIds, id.Id) {
			result.DeviceIds = append(result.DeviceIds, id.Id)
		}
	}

	devicesById, err, code := this.ListDevices(token, url.Values{"ids": {strings.Join(edit.DeviceIds, ",")}})
	if err != nil {
		return result, err, code
	}
	for _, device := range devicesById {
		if !slices.Contains(result.DeviceLocalIds, device.LocalId) {
			result.DeviceLocalIds = append(result.DeviceLocalIds, device.LocalId)
		}
	}

	if len(result.DeviceIds) != len(result.DeviceLocalIds) {
		return result, errors.New("not all local device ids found"), http.StatusBadRequest
	}
	return result, err, code
}
