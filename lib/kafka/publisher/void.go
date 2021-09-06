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
	"github.com/SENERGY-Platform/device-manager/lib/model"
)

type Void struct{}

var VoidPublisherError = errors.New("try to use void publisher")

func (this Void) PublishDevice(device model.Device, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishDeviceType(device model.DeviceType, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceTypeDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishDeviceGroup(device model.DeviceGroup, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceGroupDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishProtocol(device model.Protocol, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishProtocolDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishHub(hub model.Hub, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishHubDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishConcept(concept model.Concept, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishConceptDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishCharacteristic(conceptId string, concept model.Characteristic, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishCharacteristicDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishAspect(device model.Aspect, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishAspectDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishFunction(device model.Function, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishFunctionDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishDeviceClass(device model.DeviceClass, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishDeviceClassDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishLocation(device model.Location, userID string) (err error) {
	return VoidPublisherError
}

func (this Void) PublishLocationDelete(id string, userID string) error {
	return VoidPublisherError
}

func (this Void) PublishDeleteUserRights(resource string, id string, userId string) error {
	return VoidPublisherError
}
