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
	"github.com/SENERGY-Platform/models/go/models"
	"net/http"
	"runtime/debug"
	"sort"
	"strings"
)

func (this *Controller) ReadDeviceType(token auth.Token, id string) (dt models.DeviceType, err error, code int) {
	dt, err, code = this.com.GetDeviceType(token, id)
	sort.Slice(dt.Services, func(i, j int) bool {
		return dt.Services[i].Name < dt.Services[j].Name
	})
	return dt, err, code
}

func (this *Controller) PublishDeviceTypeCreate(token auth.Token, dt models.DeviceType) (models.DeviceType, error, int) {
	if dt.Id != "" {
		return dt, errors.New("expect empty id"), http.StatusBadRequest
	}
	dt.GenerateId()
	err, code := this.com.ValidateDeviceType(token, dt)
	if err != nil {
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, token.GetUserId())
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeUpdate(token auth.Token, id string, dt models.DeviceType) (models.DeviceType, error, int) {
	if dt.Id != id {
		return dt, errors.New("id in body unequal to id in request endpoint"), http.StatusBadRequest
	}

	dt.GenerateId()

	if !token.IsAdmin() {
		err, code := this.com.PermissionCheckForDeviceType(token, id, "w")
		if err != nil {
			debug.PrintStack()
			return dt, err, code
		}
	}
	err, code := this.com.ValidateDeviceType(token, dt)
	if err != nil {
		debug.PrintStack()
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, token.GetUserId())
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeDelete(token auth.Token, id string) (error, int) {
	if err := com.PreventIdModifier(id); err != nil {
		return err, http.StatusBadRequest
	}
	exists, err, code := this.com.DevicesOfTypeExist(token, id)
	if err != nil {
		return err, code
	}
	if exists {
		return errors.New("expect no dependent devices"), http.StatusBadRequest
	}
	err, code = this.com.PermissionCheckForDeviceType(token, id, "a")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceTypeDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (this *Controller) ValidateDistinctDeviceTypeAttributes(token auth.Token, devicetype models.DeviceType, attributeKeys []string) error {
	dtAttr := map[string]string{}
	for _, attr := range devicetype.Attributes {
		dtAttr[strings.TrimSpace(attr.Key)] = strings.TrimSpace(attr.Value)
	}

	and := []com.Selection{}
	for _, attrKey := range attributeKeys {
		attrKey = strings.TrimSpace(attrKey)
		if value, ok := dtAttr[attrKey]; ok {
			and = append(and, com.Selection{
				Condition: com.ConditionConfig{
					Feature:   "features.attributes.key",
					Operation: com.QueryEqualOperation,
					Value:     attrKey,
				},
			})
			if value != "" {
				and = append(and, com.Selection{
					Condition: com.ConditionConfig{
						Feature:   "features.attributes.value",
						Operation: com.QueryEqualOperation,
						Value:     strings.TrimSpace(value),
					},
				})
			}
		} else {
			return errors.New("distinct attribute not in device-type attributes found")
		}
	}

	var filter com.Selection
	if len(and) == 1 {
		filter = and[0]
	} else {
		filter = com.Selection{
			And: and,
		}
	}

	searchResult := []models.DeviceType{}
	err, _ := this.com.QueryPermissionsSearch(token.Jwt(), com.QueryMessage{
		Resource: this.config.DeviceTypeTopic,
		Find: &com.QueryFind{
			QueryListCommons: com.QueryListCommons{
				Limit:  1000,
				Offset: 0,
				Rights: "r",
			},
			Filter: &filter,
		},
	}, &searchResult)
	if err != nil {
		return err
	}
	for _, element := range searchResult {
		eAttr := map[string]string{}
		for _, attr := range element.Attributes {
			eAttr[attr.Key] = strings.TrimSpace(attr.Value)
		}
		if element.Id != devicetype.Id && attributesMatch(dtAttr, eAttr, attributeKeys) {
			return errors.New("find matching distinct attributes in " + element.Id)
		}
	}
	return nil
}

func attributesMatch(a, b map[string]string, attributes []string) bool {
	for _, attrKey := range attributes {
		attrKey = strings.TrimSpace(attrKey)
		aValue, ok := a[attrKey]
		if !ok {
			return false
		}
		bValue, ok := b[attrKey]
		if !ok {
			return false
		}
		if strings.TrimSpace(aValue) != strings.TrimSpace(bValue) {
			return false
		}
	}
	return true
}
