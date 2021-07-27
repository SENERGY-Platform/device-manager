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
	"github.com/SENERGY-Platform/device-manager/lib/model"
)

type Controller interface {
	ReadDeviceGroup(token string, id string) (device model.DeviceGroup, err error, code int)
	PublishDeviceGroupCreate(token string, dt model.DeviceGroup) (result model.DeviceGroup, err error, code int)
	PublishDeviceGroupUpdate(token string, id string, device model.DeviceGroup) (result model.DeviceGroup, err error, code int)
	PublishDeviceGroupDelete(token string, id string) (err error, code int)

	ReadDeviceType(token string, id string) (device model.DeviceType, err error, code int)
	PublishDeviceTypeCreate(token string, dt model.DeviceType) (result model.DeviceType, err error, code int)
	PublishDeviceTypeUpdate(token string, id string, device model.DeviceType) (result model.DeviceType, err error, code int)
	PublishDeviceTypeDelete(token string, id string) (err error, code int)

	ReadDevice(token string, id string) (device model.Device, err error, code int)
	PublishDeviceCreate(token string, dt model.Device) (result model.Device, err error, code int)
	PublishDeviceUpdate(token string, id string, device model.Device) (result model.Device, err error, code int)
	PublishDeviceDelete(token string, id string) (err error, code int)

	ReadHub(token string, id string) (hub model.Hub, err error, code int)
	PublishHubCreate(token string, dt model.Hub) (result model.Hub, err error, code int)
	PublishHubUpdate(token string, id string, hub model.Hub) (result model.Hub, err error, code int)
	PublishHubDelete(token string, id string) (err error, code int)

	ReadProtocol(token string, id string) (device model.Protocol, err error, code int)
	PublishProtocolCreate(token string, dt model.Protocol) (result model.Protocol, err error, code int)
	PublishProtocolUpdate(token string, id string, device model.Protocol) (result model.Protocol, err error, code int)
	PublishProtocolDelete(token string, id string) (err error, code int)

	PublishConceptCreate(token string, concept model.Concept) (result model.Concept, err error, code int)
	PublishConceptUpdate(token string, id string, concept model.Concept) (result model.Concept, err error, code int)
	PublishConceptDelete(token string, id string) (err error, code int)

	PublishCharacteristicCreate(token string, conceptId string, characteristic model.Characteristic) (result model.Characteristic, err error, code int)
	PublishCharacteristicUpdate(token string, conceptId string, characteristicId string, characteristic model.Characteristic) (result model.Characteristic, err error, code int)
	PublishCharacteristicDelete(token string, id string) (err error, code int)

	DeviceLocalIdToId(token string, localId string) (id string, err error, errCode int)

	ReadAspect(token string, id string) (device model.Aspect, err error, code int)
	PublishAspectCreate(token string, dt model.Aspect) (result model.Aspect, err error, code int)
	PublishAspectUpdate(token string, id string, device model.Aspect) (result model.Aspect, err error, code int)
	PublishAspectDelete(token string, id string) (err error, code int)

	ReadFunction(token string, id string) (device model.Function, err error, code int)
	PublishFunctionCreate(token string, dt model.Function) (result model.Function, err error, code int)
	PublishFunctionUpdate(token string, id string, device model.Function) (result model.Function, err error, code int)
	PublishFunctionDelete(token string, id string) (err error, code int)

	ReadDeviceClass(token string, id string) (device model.DeviceClass, err error, code int)
	PublishDeviceClassCreate(token string, dt model.DeviceClass) (result model.DeviceClass, err error, code int)
	PublishDeviceClassUpdate(token string, id string, device model.DeviceClass) (result model.DeviceClass, err error, code int)
	PublishDeviceClassDelete(token string, id string) (err error, code int)

	ReadLocation(token string, id string) (device model.Location, err error, code int)
	PublishLocationCreate(token string, dt model.Location) (result model.Location, err error, code int)
	PublishLocationUpdate(token string, id string, device model.Location) (result model.Location, err error, code int)
	PublishLocationDelete(token string, id string) (err error, code int)
}
