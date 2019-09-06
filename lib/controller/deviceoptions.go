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
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

func (this *Controller) FindDeviceOptions(jwt jwt_http_router.Jwt, descriptions []model.DeviceDescription) (result []model.DeviceOption, err error, code int) {
	deviceTypes, err, code := this.com.GetDeviceTypeFromDescriptions(jwt, descriptions)
	if err != nil {
		return result, err, code
	}
	for _, dt := range deviceTypes {
		devices, err, code := this.com.GetDevicesOfType(jwt, dt.Id)
		if err != nil {
			return result, err, code
		}
		services := []model.Service{}
		serviceIndex := map[string]model.Service{}
		for _, service := range dt.Services {
			for _, desc := range descriptions {
				for _, function := range service.Functions {
					if function.Id == desc.Function.Id {
						if desc.Aspect == nil {
							serviceIndex[service.Id] = service
						} else {
							for _, aspect := range service.Aspects {
								if aspect.Id == desc.Aspect.Id {
									serviceIndex[service.Id] = service
								}
							}
						}
					}
				}
			}
		}
		for _, service := range serviceIndex {
			services = append(services, service)
		}
		for _, device := range devices {
			result = append(result, model.DeviceOption{
				Device:         device,
				ServiceOptions: services,
			})
		}
	}
	return result, nil, 200
}
