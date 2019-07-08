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

package mock

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
)

const DtTopic = "device-type"

type Publisher struct {
	listener map[string][]func(msg []byte)
}

func NewPublisher() *Publisher {
	return &Publisher{listener: map[string][]func(msg []byte){}}
}

func (this *Publisher) PublishDeviceType(device model.DeviceType, userId string) (err error) {
	cmd := publisher.DeviceTypeCommand{Command: "PUT", Id: device.Id, DeviceType: device, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DtTopic, message)
}

func (this *Publisher) PublishDeviceDelete(id string, userId string) error {
	cmd := publisher.DeviceTypeCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DtTopic, message)
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
