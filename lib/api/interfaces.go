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
}
