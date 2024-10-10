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
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
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

func (this *Controller) PublishDeviceTypeCreate(token auth.Token, dt models.DeviceType, options model.DeviceTypeUpdateOptions) (models.DeviceType, error, int) {
	if dt.Id != "" {
		return dt, errors.New("expect empty id"), http.StatusBadRequest
	}
	dt.GenerateId()
	err, code := this.com.ValidateDeviceType(token, dt)
	if err != nil {
		return dt, err, code
	}

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceTypeTopic,
		ResourceId:   dt.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishDeviceType(dt, token.GetUserId())
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}

	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeUpdate(token auth.Token, id string, dt models.DeviceType, options model.DeviceTypeUpdateOptions) (models.DeviceType, error, int) {
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

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceTypeTopic,
		ResourceId:   dt.Id,
		Command:      "PUT",
	})

	err = this.publisher.PublishDeviceType(dt, token.GetUserId())
	if err != nil {
		debug.PrintStack()
		return dt, err, http.StatusInternalServerError
	}

	err = wait()
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}

	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeDelete(token auth.Token, id string, options model.DeviceTypeDeleteOptions) (error, int) {
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

	wait := this.optionalWait(options.Wait, donewait.DoneMsg{
		ResourceKind: this.config.DeviceTypeTopic,
		ResourceId:   id,
		Command:      "DELETE",
	})

	err = this.publisher.PublishDeviceTypeDelete(id, token.GetUserId())
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = wait()
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

	options := client.DeviceTypeListOptions{
		Limit:           9999,
		AttributeKeys:   nil,
		AttributeValues: nil,
	}

	for _, attrKey := range attributeKeys {
		attrKey = strings.TrimSpace(attrKey)
		if value, ok := dtAttr[attrKey]; ok {
			options.AttributeKeys = append(options.AttributeKeys, attrKey)
			options.AttributeValues = append(options.AttributeValues, strings.TrimSpace(value))
		} else {
			return errors.New("distinct attribute not in device-type attributes found")
		}
	}
	list, err, _ := this.com.ListDeviceTypes(token.Jwt(), options)
	if err != nil {
		return err
	}
	for _, element := range list {
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
