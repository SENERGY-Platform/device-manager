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
)

type Controller interface {
	ReadDeviceGroup(token auth.Token, id string) (device model.DeviceGroup, err error, code int)
	PublishDeviceGroupCreate(token auth.Token, dt model.DeviceGroup) (result model.DeviceGroup, err error, code int)
	PublishDeviceGroupUpdate(token auth.Token, id string, device model.DeviceGroup) (result model.DeviceGroup, err error, code int)
	PublishDeviceGroupDelete(token auth.Token, id string) (err error, code int)

	ReadDeviceType(token auth.Token, id string) (device model.DeviceType, err error, code int)
	PublishDeviceTypeCreate(token auth.Token, dt model.DeviceType) (result model.DeviceType, err error, code int)
	PublishDeviceTypeUpdate(token auth.Token, id string, device model.DeviceType) (result model.DeviceType, err error, code int)
	PublishDeviceTypeDelete(token auth.Token, id string) (err error, code int)

	ReadDevice(token auth.Token, id string) (device model.Device, err error, code int)
	PublishDeviceCreate(token auth.Token, dt model.Device) (result model.Device, err error, code int)
	PublishDeviceUpdate(token auth.Token, id string, device model.Device) (result model.Device, err error, code int)
	PublishDeviceDelete(token auth.Token, id string) (err error, code int)

	ReadHub(token auth.Token, id string) (hub model.Hub, err error, code int)
	PublishHubCreate(token auth.Token, dt model.Hub) (result model.Hub, err error, code int)
	PublishHubUpdate(token auth.Token, id string, hub model.Hub) (result model.Hub, err error, code int)
	PublishHubDelete(token auth.Token, id string) (err error, code int)

	ReadProtocol(token auth.Token, id string) (device model.Protocol, err error, code int)
	PublishProtocolCreate(token auth.Token, dt model.Protocol) (result model.Protocol, err error, code int)
	PublishProtocolUpdate(token auth.Token, id string, device model.Protocol) (result model.Protocol, err error, code int)
	PublishProtocolDelete(token auth.Token, id string) (err error, code int)

	PublishConceptCreate(token auth.Token, concept model.Concept) (result model.Concept, err error, code int)
	PublishConceptUpdate(token auth.Token, id string, concept model.Concept) (result model.Concept, err error, code int)
	PublishConceptDelete(token auth.Token, id string) (err error, code int)

	PublishCharacteristicCreate(token auth.Token, conceptId string, characteristic model.Characteristic) (result model.Characteristic, err error, code int)
	PublishCharacteristicUpdate(token auth.Token, conceptId string, characteristicId string, characteristic model.Characteristic) (result model.Characteristic, err error, code int)
	PublishCharacteristicDelete(token auth.Token, id string) (err error, code int)

	DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, errCode int)

	ReadAspect(token auth.Token, id string) (device model.Aspect, err error, code int)
	PublishAspectCreate(token auth.Token, dt model.Aspect) (result model.Aspect, err error, code int)
	PublishAspectUpdate(token auth.Token, id string, device model.Aspect) (result model.Aspect, err error, code int)
	PublishAspectDelete(token auth.Token, id string) (err error, code int)

	ReadFunction(token auth.Token, id string) (device model.Function, err error, code int)
	PublishFunctionCreate(token auth.Token, dt model.Function) (result model.Function, err error, code int)
	PublishFunctionUpdate(token auth.Token, id string, device model.Function) (result model.Function, err error, code int)
	PublishFunctionDelete(token auth.Token, id string) (err error, code int)

	ReadDeviceClass(token auth.Token, id string) (device model.DeviceClass, err error, code int)
	PublishDeviceClassCreate(token auth.Token, dt model.DeviceClass) (result model.DeviceClass, err error, code int)
	PublishDeviceClassUpdate(token auth.Token, id string, device model.DeviceClass) (result model.DeviceClass, err error, code int)
	PublishDeviceClassDelete(token auth.Token, id string) (err error, code int)

	ReadLocation(token auth.Token, id string) (device model.Location, err error, code int)
	PublishLocationCreate(token auth.Token, dt model.Location) (result model.Location, err error, code int)
	PublishLocationUpdate(token auth.Token, id string, device model.Location) (result model.Location, err error, code int)
	PublishLocationDelete(token auth.Token, id string) (err error, code int)
}
