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
	"github.com/SENERGY-Platform/models/go/models"
	permmodel "github.com/SENERGY-Platform/permission-search/lib/model"
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
	PublishDevice(device models.Device, userID string) (err error)
	PublishDeviceDelete(id string, userID string) error

	PublishDeviceType(device models.DeviceType, userID string) (err error)
	PublishDeviceTypeDelete(id string, userID string) error

	PublishDeviceGroup(device models.DeviceGroup, userID string) (err error)
	PublishDeviceGroupDelete(id string, userID string) error

	PublishProtocol(device models.Protocol, userID string) (err error)
	PublishProtocolDelete(id string, userID string) error

	PublishHub(hub models.Hub, userID string) (err error)
	PublishHubDelete(id string, userID string) error

	PublishConcept(concept models.Concept, userID string) (err error)
	PublishConceptDelete(id string, userID string) error

	PublishCharacteristic(characteristic models.Characteristic, userID string) (err error)
	PublishCharacteristicDelete(id string, userID string) error

	PublishAspect(device models.Aspect, userID string) (err error)
	PublishAspectDelete(id string, userID string) error

	PublishFunction(device models.Function, userID string) (err error)
	PublishFunctionDelete(id string, userID string) error

	PublishDeviceClass(device models.DeviceClass, userID string) (err error)
	PublishDeviceClassDelete(id string, userID string) error

	PublishLocation(device models.Location, userID string) (err error)
	PublishLocationDelete(id string, userID string) error

	PublishRights(kind string, id string, element permmodel.ResourceRightsBase) error
}

type Com interface {
	ResourcesEffectedByUserDelete(token auth.Token, resource string) (deleteResourceIds []string, deleteUserFromResourceIds []string, err error)
	GetPermissions(token auth.Token, kind string, id string) (permmodel.ResourceRights, error)

	GetTechnicalDeviceGroup(token auth.Token, id string) (dt models.DeviceGroup, err error, code int)
	ValidateDeviceGroup(token auth.Token, dt models.DeviceGroup) (err error, code int)
	PermissionCheckForDeviceGroup(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetDeviceType(token auth.Token, id string) (dt models.DeviceType, err error, code int)
	ValidateDeviceType(token auth.Token, dt models.DeviceType) (err error, code int)
	PermissionCheckForDeviceType(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetProtocol(token auth.Token, id string) (models.Protocol, error, int)
	ValidateProtocol(token auth.Token, protocol models.Protocol) (err error, code int)

	GetDevice(token auth.Token, id string) (models.Device, error, int) //uses internal admin jwt
	ValidateDevice(token auth.Token, device models.Device) (err error, code int)
	PermissionCheckForDevice(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"
	PermissionCheckForDeviceList(token auth.Token, ids []string, rights string) (result map[string]bool, err error, code int)

	GetHub(token auth.Token, id string) (models.Hub, error, int) //uses internal admin jwt
	ValidateHub(token auth.Token, hub models.Hub) (err error, code int)
	PermissionCheckForHub(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	GetConcept(token auth.Token, id string) (models.Concept, error, int)
	ValidateConcept(token auth.Token, concept models.Concept) (err error, code int)
	PermissionCheckForConcept(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	ValidateCharacteristic(token auth.Token, concept models.Characteristic) (err error, code int)
	PermissionCheckForCharacteristic(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"
	GetCharacteristic(token auth.Token, id string) (concept models.Characteristic, err error, code int)

	DevicesOfTypeExist(token auth.Token, deviceTypeId string) (result bool, err error, code int)

	DeviceLocalIdToId(token auth.Token, localId string) (id string, err error, code int)

	GetAspect(token auth.Token, id string) (models.Aspect, error, int)
	ValidateAspect(token auth.Token, aspect models.Aspect) (err error, code int)

	GetFunction(token auth.Token, id string) (models.Function, error, int)
	ValidateFunction(token auth.Token, function models.Function) (err error, code int)

	GetDeviceClass(token auth.Token, id string) (models.DeviceClass, error, int)
	ValidateDeviceClass(token auth.Token, deviceClass models.DeviceClass) (err error, code int)

	GetLocation(token auth.Token, id string) (models.Location, error, int)
	ValidateLocation(token auth.Token, Location models.Location) (err error, code int)
	PermissionCheckForLocation(token auth.Token, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"

	ValidateAspectDelete(token auth.Token, id string) (err error, code int)
	ValidateCharacteristicDelete(token auth.Token, id string) (err error, code int)
	ValidateConceptDelete(token auth.Token, id string) (err error, code int)
	ValidateDeviceClassDelete(token auth.Token, id string) (err error, code int)
	ValidateFunctionDelete(token auth.Token, id string) (err error, code int)

	QueryPermissionsSearch(token string, query com.QueryMessage, result interface{}) (err error, code int)
}

func (this *Controller) GetPublisher() Publisher {
	return this.publisher
}
