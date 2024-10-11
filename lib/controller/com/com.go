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
	"github.com/SENERGY-Platform/device-manager/lib/config"
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/permissions-v2/pkg/client"
)

type Com struct {
	config  config.Config
	perm    client.Client
	devices devicerepo.Interface
}

func New(config config.Config) *Com {
	return &Com{
		config:  config,
		perm:    client.New(config.PermissionsV2Url),
		devices: devicerepo.NewClient(config.DeviceRepoUrl),
	}
}
