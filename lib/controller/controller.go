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
	"context"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/listener"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/publisher"
	"github.com/SENERGY-Platform/device-manager/lib/model"
)

type Controller struct {
	publisher Publisher
	com       Com
	config    config.Config
}

func New(basectx context.Context, conf config.Config) (ctrl *Controller, err error) {
	ctx, cancel := context.WithCancel(basectx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	var publ Publisher
	if conf.EditForward == "" || conf.EditForward == "-" {
		publ, err = publisher.New(conf, ctx)
		if err != nil {
			return &Controller{}, err
		}
	} else {
		publ = publisher.Void{}
	}
	ctrl = &Controller{com: com.New(conf), publisher: publ, config: conf}
	if conf.EditForward == "" || conf.EditForward == "-" {
		err = listener.Start(ctx, conf, ctrl)
	}
	return
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

	PublishDeleteUserRights(resource string, id string, userId string) error
}

type Com interface {
	ResourcesEffectedByUserDelete(token auth.Token, resource string) (deleteResourceIds []string, deleteUserFromResourceIds []string, err error)

	GetTechnicalDeviceGroup(token auth.Token, id string) (dt model.DeviceGroup, err error, code int)
	ValidateDeviceGroup(token auth.Token, dt model.DeviceGroup) (err error, code int)
	PermissionCheckForDeviceGroup(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetTechnicalDeviceType(token auth.Token, id string) (dt model.DeviceType, err error, code int)
	GetSemanticDeviceType(token auth.Token, id string) (dt model.DeviceType, err error, code int)
	ValidateDeviceType(token auth.Token, dt model.DeviceType) (err error, code int)
	PermissionCheckForDeviceType(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetProtocol(token auth.Token, id string) (model.Protocol, error, int)
	ValidateProtocol(token auth.Token, protocol model.Protocol) (err error, code int)

	GetDevice(token auth.Token, id string) (model.Device, error, int) //uses internal admin jwt
	ValidateDevice(token auth.Token, device model.Device) (err error, code int)
	PermissionCheckForDevice(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetHub(token auth.Token, id string) (model.Hub, error, int) //uses internal admin jwt
	ValidateHub(token auth.Token, hub model.Hub) (err error, code int)
	PermissionCheckForHub(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetConcept(token auth.Token, id string) (model.Concept, error, int)
	ValidateConcept(token auth.Token, concept model.Concept) (err error, code int)
	PermissionCheckForConcept(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	ValidateCharacteristic(token auth.Token, concept model.Characteristic) (err error, code int)
	PermissionCheckForCharacteristic(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	DevicesOfTypeExist(token auth.Token, deviceTypeId string) (result bool, err error, code int)

	DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, code int)

	GetAspect(token auth.Token, id string) (model.Aspect, error, int)
	ValidateAspect(token auth.Token, aspect model.Aspect) (err error, code int)

	GetFunction(token auth.Token, id string) (model.Function, error, int)
	ValidateFunction(token auth.Token, function model.Function) (err error, code int)

	GetDeviceClass(token auth.Token, id string) (model.DeviceClass, error, int)
	ValidateDeviceClass(token auth.Token, deviceClass model.DeviceClass) (err error, code int)

	GetLocation(token auth.Token, id string) (model.Location, error, int)
	ValidateLocation(token auth.Token, Location model.Location) (err error, code int)
	PermissionCheckForLocation(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"
}
