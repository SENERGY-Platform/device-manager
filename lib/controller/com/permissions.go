/*
 * Copyright 2024 InfAI (CC SES)
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

package com

import (
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/permissions-v2/pkg/model"
)

func (this *Com) SetPermission(token string, topicId string, id string, permissions model.ResourcePermissions) (result model.ResourcePermissions, err error, code int) {
	return this.perm.SetPermission(token, topicId, id, permissions)
}

func (this *Com) ListDeviceTypes(token string, options devicerepo.DeviceTypeListOptions) (result []models.DeviceType, err error, code int) {
	result, _, err, code = this.devices.ListDeviceTypesV3(token, options)
	return
}

func (this *Com) ListDevices(token string, options devicerepo.DeviceListOptions) (result []models.Device, err error, code int) {
	return this.devices.ListDevices(token, options)
}
