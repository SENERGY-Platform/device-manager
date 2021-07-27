/*
 * Copyright 2021 InfAI (CC SES)
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
)

func (this *Com) GetLocation(token string, id string) (Location model.Location, err error, code int) {
	err, code = getResourceFromService(token, this.config.SemanticRepoUrl+"/locations", id, &Location)
	return
}

func (this *Com) ValidateLocation(token string, Location model.Location) (err error, code int) {
	return validateResource(token, []string{
		this.config.SemanticRepoUrl + "/locations?dry-run=true",
	}, Location)
}