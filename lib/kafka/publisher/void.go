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

package publisher

import (
	"errors"
	"github.com/SENERGY-Platform/models/go/models"
	permmodel "github.com/SENERGY-Platform/permission-search/lib/model"
)

type Void struct{}

var VoidPublisherError = errors.New("try to use void publisher")

func (this Void) PublishDevice(device models.Device, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishDeviceType(device models.DeviceType, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceTypeDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishDeviceGroup(device models.DeviceGroup, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceGroupDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishProtocol(device models.Protocol, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishProtocolDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishHub(hub models.Hub, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishHubDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishConcept(concept models.Concept, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishConceptDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishCharacteristic(characteristic models.Characteristic, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishCharacteristicDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishAspect(device models.Aspect, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishAspectDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishFunction(device models.Function, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishFunctionDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishDeviceClass(device models.DeviceClass, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceClassDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishLocation(device models.Location, userID string, strictWaitBeforeDone bool) (err error) {
	return VoidPublisherError
}

func (this Void) PublishLocationDelete(id string, userID string, strictWaitBeforeDone bool) error {
	return VoidPublisherError
}

func (this Void) PublishRights(kind string, id string, element permmodel.ResourceRightsBase) error {
	return VoidPublisherError
}
