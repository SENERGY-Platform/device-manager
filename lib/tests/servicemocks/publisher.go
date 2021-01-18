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

package servicemocks

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
)

var DtTopic = "devicetype"
var ProtocolTopic = "protocol"
var DeviceTopic = "device"
var DeviceGroupTopic = "devicegroups"
var HubTopic = "hub"
var ConceptTopic = "concepts"
var CharacteristicTopic = "characteristics"
var AspectTopic = "aspcet"
var FunctionTopic = "function"
var DeviceClassTopic = "function"
var LocationTopic = "locations"

type Publisher struct {
	listener map[string][]func(msg []byte)
}

func NewPublisher() *Publisher {
	return &Publisher{listener: map[string][]func(msg []byte){}}
}

func (this *Publisher) send(topic string, msg []byte) error {
	for _, listener := range this.listener[topic] {
		go listener(msg)
	}
	return nil
}

func (this *Publisher) Subscribe(topic string, f func(msg []byte)) {
	this.listener[topic] = append(this.listener[topic], f)
}

func (this *Publisher) PublishDevice(device model.Device, userId string) (err error) {
	cmd := publisher.DeviceCommand{Command: "PUT", Id: device.Id, Device: device, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DeviceTopic, message)
}

func (this *Publisher) PublishDeviceDelete(id string, userId string) error {
	cmd := publisher.DeviceCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DeviceTopic, message)
}

func (this *Publisher) PublishHub(hub model.Hub, userId string) (err error) {
	cmd := publisher.HubCommand{Command: "PUT", Id: hub.Id, Hub: hub, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(HubTopic, message)
}

func (this *Publisher) PublishHubDelete(id string, userId string) error {
	cmd := publisher.HubCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(HubTopic, message)
}

func (this *Publisher) PublishDeviceType(device model.DeviceType, userId string) (err error) {
	cmd := publisher.DeviceTypeCommand{Command: "PUT", Id: device.Id, DeviceType: device, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DtTopic, message)
}

func (this *Publisher) PublishDeviceTypeDelete(id string, userId string) error {
	cmd := publisher.DeviceTypeCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DtTopic, message)
}

func (this *Publisher) PublishDeviceGroup(group model.DeviceGroup, userID string) (err error) {
	cmd := publisher.DeviceGroupCommand{Command: "PUT", Id: group.Id, DeviceGroup: group, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DeviceGroupTopic, message)
}

func (this *Publisher) PublishDeviceGroupDelete(id string, userID string) error {
	cmd := publisher.DeviceGroupCommand{Command: "DELETE", Id: id, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DeviceGroupTopic, message)
}

func (this *Publisher) PublishProtocol(protocol model.Protocol, userId string) (err error) {
	cmd := publisher.ProtocolCommand{Command: "PUT", Id: protocol.Id, Protocol: protocol, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(ProtocolTopic, message)
}

func (this *Publisher) PublishProtocolDelete(id string, userId string) error {
	cmd := publisher.ProtocolCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(ProtocolTopic, message)
}

func (this *Publisher) PublishConcept(concept model.Concept, userID string) (err error) {
	cmd := publisher.ConceptCommand{Command: "PUT", Id: concept.Id, Concept: concept, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(ConceptTopic, message)
}

func (this *Publisher) PublishConceptDelete(id string, userID string) error {
	cmd := publisher.ProtocolCommand{Command: "DELETE", Id: id, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(ConceptTopic, message)
}

func (this *Publisher) PublishCharacteristic(conceptId string, characteristic model.Characteristic, userID string) (err error) {
	cmd := publisher.CharacteristicCommand{Command: "PUT", ConceptId: conceptId, Id: characteristic.Id, Characteristic: characteristic, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(CharacteristicTopic, message)
}

func (this *Publisher) PublishCharacteristicDelete(id string, userID string) error {
	cmd := publisher.CharacteristicCommand{Command: "DELETE", Id: id, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(CharacteristicTopic, message)
}

func (this *Publisher) PublishAspect(aspect model.Aspect, userID string) (err error) {
	cmd := publisher.AspectCommand{Command: "PUT", Id: aspect.Id, Aspect: aspect, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(AspectTopic, message)
}

func (this *Publisher) PublishAspectDelete(id string, userId string) error {
	cmd := publisher.AspectCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(AspectTopic, message)
}

func (this *Publisher) PublishFunction(function model.Function, userID string) (err error) {
	cmd := publisher.FunctionCommand{Command: "PUT", Id: function.Id, Function: function, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(FunctionTopic, message)
}

func (this *Publisher) PublishFunctionDelete(id string, userId string) error {
	cmd := publisher.FunctionCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(FunctionTopic, message)
}

func (this *Publisher) PublishDeviceClass(deviceClass model.DeviceClass, userID string) (err error) {
	cmd := publisher.DeviceClassCommand{Command: "PUT", Id: deviceClass.Id, DeviceClass: deviceClass, Owner: userID}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DeviceClassTopic, message)
}

func (this *Publisher) PublishDeviceClassDelete(id string, userId string) error {
	cmd := publisher.DeviceClassCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DeviceClassTopic, message)
}

func (this *Publisher) PublishLocation(location model.Location, userId string) (err error) {
	cmd := publisher.LocationCommand{Command: "PUT", Id: location.Id, Location: location, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(LocationTopic, message)
}

func (this *Publisher) PublishLocationDelete(id string, userId string) error {
	cmd := publisher.LocationCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(LocationTopic, message)
}
