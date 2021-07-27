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
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
)

type Controller struct {
	publisher Publisher
	com       Com
}

func New(conf config.Config) (*Controller, error) {
	publ, err := publisher.New(conf)
	if err != nil {
		return &Controller{}, err
	}
	return &Controller{com: com.New(conf), publisher: publ}, nil
}

func NewWithPublisher(conf config.Config, publisher Publisher) (*Controller, error) {
	return &Controller{com: com.New(conf), publisher: publisher}, nil
}

type Publisher interface {
	PublishDevice(device model.Device, userID string) (err error)
	PublishDeviceDelete(id string, userID string) error

	PublishDeviceType(device model.DeviceType, userID string) (err error)
	PublishDeviceTypeDelete(id string, userID string) error

	PublishDeviceGroup(device model.DeviceGroup, userID string) (err error)
	PublishDeviceGroupDelete(id string, userID string) error

	PublishProtocol(device model.Protocol, userID string) (err error)
	PublishProtocolDelete(id string, userID string) error

	PublishHub(hub model.Hub, userID string) (err error)
	PublishHubDelete(id string, userID string) error

	PublishConcept(concept model.Concept, userID string) (err error)
	PublishConceptDelete(id string, userID string) error

	PublishCharacteristic(conceptId string, concept model.Characteristic, userID string) (err error)
	PublishCharacteristicDelete(id string, userID string) error

	PublishAspect(device model.Aspect, userID string) (err error)
	PublishAspectDelete(id string, userID string) error

	PublishFunction(device model.Function, userID string) (err error)
	PublishFunctionDelete(id string, userID string) error

	PublishDeviceClass(device model.DeviceClass, userID string) (err error)
	PublishDeviceClassDelete(id string, userID string) error

	PublishLocation(device model.Location, userID string) (err error)
	PublishLocationDelete(id string, userID string) error
}

type Com interface {
	GetTechnicalDeviceGroup(token string, id string) (dt model.DeviceGroup, err error, code int)
	ValidateDeviceGroup(token string, dt model.DeviceGroup) (err error, code int)
	PermissionCheckForDeviceGroup(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetTechnicalDeviceType(token string, id string) (dt model.DeviceType, err error, code int)
	GetSemanticDeviceType(token string, id string) (dt model.DeviceType, err error, code int)
	ValidateDeviceType(token string, dt model.DeviceType) (err error, code int)
	PermissionCheckForDeviceType(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetProtocol(token string, id string) (model.Protocol, error, int)
	ValidateProtocol(token string, protocol model.Protocol) (err error, code int)

	GetDevice(token string, id string) (model.Device, error, int) //uses internal admin jwt
	ValidateDevice(token string, device model.Device) (err error, code int)
	PermissionCheckForDevice(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetHub(token string, id string) (model.Hub, error, int) //uses internal admin jwt
	ValidateHub(token string, hub model.Hub) (err error, code int)
	PermissionCheckForHub(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	ValidateConcept(token string, concept model.Concept) (err error, code int)
	PermissionCheckForConcept(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	ValidateCharacteristic(token string, concept model.Characteristic) (err error, code int)
	PermissionCheckForCharacteristic(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	DevicesOfTypeExist(token string, deviceTypeId string) (result bool, err error, code int)

	DeviceLocalIdToId(token string, localId string) (id string, err error, code int)

	GetAspect(token string, id string) (model.Aspect, error, int)
	ValidateAspect(token string, aspect model.Aspect) (err error, code int)

	GetFunction(token string, id string) (model.Function, error, int)
	ValidateFunction(token string, function model.Function) (err error, code int)

	GetDeviceClass(token string, id string) (model.DeviceClass, error, int)
	ValidateDeviceClass(token string, deviceClass model.DeviceClass) (err error, code int)

	GetLocation(token string, id string) (model.Location, error, int)
	ValidateLocation(token string, Location model.Location) (err error, code int)
	PermissionCheckForLocation(token string, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"
}
