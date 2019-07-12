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

package com

import (
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

//expects previous permission check and use own admin jwt to access hub
func (this *Com) GetHub(jwt jwt_http_router.Jwt, id string) (hub model.Hub, err error, code int) {
	err, code = getResourceFromService(jwt, this.config.DeviceRepoUrl+"/hubs", id, &hub)
	return
}

func (this *Com) ValidateHub(jwt jwt_http_router.Jwt, hub model.Hub) (err error, code int) {
	return validateResource(jwt, []string{
		this.config.DeviceRepoUrl + "/hubs?dry-run=true",
	}, hub)
}
