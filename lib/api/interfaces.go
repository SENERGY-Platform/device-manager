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

package api

import (
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
)

type Controller interface {
	ReadDeviceGroup(token auth.Token, id string) (device models.DeviceGroup, err error, code int)
	PublishDeviceGroupCreate(token auth.Token, dt models.DeviceGroup) (result models.DeviceGroup, err error, code int)
	PublishDeviceGroupUpdate(token auth.Token, id string, device models.DeviceGroup) (result models.DeviceGroup, err error, code int)
	PublishDeviceGroupDelete(token auth.Token, id string) (err error, code int)

	ReadDeviceType(token auth.Token, id string) (device models.DeviceType, err error, code int)
	PublishDeviceTypeCreate(token auth.Token, dt models.DeviceType) (result models.DeviceType, err error, code int)
	PublishDeviceTypeUpdate(token auth.Token, id string, device models.DeviceType) (result models.DeviceType, err error, code int)
	PublishDeviceTypeDelete(token auth.Token, id string) (err error, code int)

	ReadDevice(token auth.Token, id string) (device models.Device, err error, code int)
	PublishDeviceCreate(token auth.Token, dt models.Device) (result models.Device, err error, code int)
	PublishDeviceUpdate(token auth.Token, id string, device models.Device, options model.DeviceUpdateOptions) (result models.Device, err error, code int)
	PublishDeviceDelete(token auth.Token, id string) (err error, code int)

	ReadHub(token auth.Token, id string) (hub models.Hub, err error, code int)
	PublishHubCreate(token auth.Token, dt models.HubEdit) (result models.Hub, err error, code int)
	PublishHubUpdate(token auth.Token, id string, userId string, hub models.HubEdit) (result models.Hub, err error, code int)
	PublishHubDelete(token auth.Token, id string) (err error, code int)

	ReadProtocol(token auth.Token, id string) (device models.Protocol, err error, code int)
	PublishProtocolCreate(token auth.Token, dt models.Protocol) (result models.Protocol, err error, code int)
	PublishProtocolUpdate(token auth.Token, id string, device models.Protocol) (result models.Protocol, err error, code int)
	PublishProtocolDelete(token auth.Token, id string) (err error, code int)

	ReadConcept(token auth.Token, id string) (device models.Concept, err error, code int)
	PublishConceptCreate(token auth.Token, concept models.Concept) (result models.Concept, err error, code int)
	PublishConceptUpdate(token auth.Token, id string, concept models.Concept) (result models.Concept, err error, code int)
	PublishConceptDelete(token auth.Token, id string) (err error, code int)

	PublishCharacteristicCreate(token auth.Token, characteristic models.Characteristic) (result models.Characteristic, err error, code int)
	PublishCharacteristicUpdate(token auth.Token, characteristicId string, characteristic models.Characteristic) (result models.Characteristic, err error, code int)
	PublishCharacteristicDelete(token auth.Token, id string) (err error, code int)
	ReadCharacteristic(token auth.Token, id string) (result models.Characteristic, err error, code int)

	DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, errCode int)

	ReadAspect(token auth.Token, id string) (device models.Aspect, err error, code int)
	PublishAspectCreate(token auth.Token, dt models.Aspect) (result models.Aspect, err error, code int)
	PublishAspectUpdate(token auth.Token, id string, device models.Aspect) (result models.Aspect, err error, code int)
	PublishAspectDelete(token auth.Token, id string) (err error, code int)

	ReadFunction(token auth.Token, id string) (device models.Function, err error, code int)
	PublishFunctionCreate(token auth.Token, dt models.Function) (result models.Function, err error, code int)
	PublishFunctionUpdate(token auth.Token, id string, device models.Function) (result models.Function, err error, code int)
	PublishFunctionDelete(token auth.Token, id string) (err error, code int)

	ReadDeviceClass(token auth.Token, id string) (device models.DeviceClass, err error, code int)
	PublishDeviceClassCreate(token auth.Token, dt models.DeviceClass) (result models.DeviceClass, err error, code int)
	PublishDeviceClassUpdate(token auth.Token, id string, device models.DeviceClass) (result models.DeviceClass, err error, code int)
	PublishDeviceClassDelete(token auth.Token, id string) (err error, code int)

	ReadLocation(token auth.Token, id string) (device models.Location, err error, code int)
	PublishLocationCreate(token auth.Token, dt models.Location) (result models.Location, err error, code int)
	PublishLocationUpdate(token auth.Token, id string, device models.Location) (result models.Location, err error, code int)
	PublishLocationDelete(token auth.Token, id string) (err error, code int)

	ValidateDistinctDeviceTypeAttributes(token auth.Token, devicetype models.DeviceType, attributeKeys []string) error
}
