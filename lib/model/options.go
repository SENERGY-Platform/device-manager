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

package model

type DeviceDescription struct {
	CharacteristicId string       `json:"characteristic_id"`
	Function         Function     `json:"function"`
	DeviceClass      *DeviceClass `json:"device_class,omitempty"`
	Aspect           *Aspect      `json:"aspect,omitempty"`
}

type DeviceOption struct {
	Device         Device    `json:"device"`
	ServiceOptions []Service `json:"service_options"`
}
