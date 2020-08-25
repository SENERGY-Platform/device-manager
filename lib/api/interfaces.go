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
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

type Controller interface {
	ReadDeviceType(jwt jwt_http_router.Jwt, id string) (device model.DeviceType, err error, code int)
	PublishDeviceTypeCreate(jwt jwt_http_router.Jwt, dt model.DeviceType) (result model.DeviceType, err error, code int)
	PublishDeviceTypeUpdate(jwt jwt_http_router.Jwt, id string, device model.DeviceType) (result model.DeviceType, err error, code int)
	PublishDeviceTypeDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	ReadDevice(jwt jwt_http_router.Jwt, id string) (device model.Device, err error, code int)
	PublishDeviceCreate(jwt jwt_http_router.Jwt, dt model.Device) (result model.Device, err error, code int)
	PublishDeviceUpdate(jwt jwt_http_router.Jwt, id string, device model.Device) (result model.Device, err error, code int)
	PublishDeviceDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	ReadHub(jwt jwt_http_router.Jwt, id string) (hub model.Hub, err error, code int)
	PublishHubCreate(jwt jwt_http_router.Jwt, dt model.Hub) (result model.Hub, err error, code int)
	PublishHubUpdate(jwt jwt_http_router.Jwt, id string, hub model.Hub) (result model.Hub, err error, code int)
	PublishHubDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	ReadProtocol(jwt jwt_http_router.Jwt, id string) (device model.Protocol, err error, code int)
	PublishProtocolCreate(jwt jwt_http_router.Jwt, dt model.Protocol) (result model.Protocol, err error, code int)
	PublishProtocolUpdate(jwt jwt_http_router.Jwt, id string, device model.Protocol) (result model.Protocol, err error, code int)
	PublishProtocolDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	PublishConceptCreate(jwt jwt_http_router.Jwt, concept model.Concept) (result model.Concept, err error, code int)
	PublishConceptUpdate(jwt jwt_http_router.Jwt, id string, concept model.Concept) (result model.Concept, err error, code int)
	PublishConceptDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	PublishCharacteristicCreate(jwt jwt_http_router.Jwt, conceptId string, characteristic model.Characteristic) (result model.Characteristic, err error, code int)
	PublishCharacteristicUpdate(jwt jwt_http_router.Jwt, conceptId string, characteristicId string, characteristic model.Characteristic) (result model.Characteristic, err error, code int)
	PublishCharacteristicDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	DeviceLocalIdToId(jwt jwt_http_router.Jwt, localId string) (id string, err error, errCode int)

	ReadAspect(jwt jwt_http_router.Jwt, id string) (device model.Aspect, err error, code int)
	PublishAspectCreate(jwt jwt_http_router.Jwt, dt model.Aspect) (result model.Aspect, err error, code int)
	PublishAspectUpdate(jwt jwt_http_router.Jwt, id string, device model.Aspect) (result model.Aspect, err error, code int)
	PublishAspectDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	ReadFunction(jwt jwt_http_router.Jwt, id string) (device model.Function, err error, code int)
	PublishFunctionCreate(jwt jwt_http_router.Jwt, dt model.Function) (result model.Function, err error, code int)
	PublishFunctionUpdate(jwt jwt_http_router.Jwt, id string, device model.Function) (result model.Function, err error, code int)
	PublishFunctionDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)

	ReadDeviceClass(jwt jwt_http_router.Jwt, id string) (device model.DeviceClass, err error, code int)
	PublishDeviceClassCreate(jwt jwt_http_router.Jwt, dt model.DeviceClass) (result model.DeviceClass, err error, code int)
	PublishDeviceClassUpdate(jwt jwt_http_router.Jwt, id string, device model.DeviceClass) (result model.DeviceClass, err error, code int)
	PublishDeviceClassDelete(jwt jwt_http_router.Jwt, id string) (err error, code int)
}
