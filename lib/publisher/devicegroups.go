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

package publisher

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/segmentio/kafka-go"
	"log"
	"runtime/debug"
	"time"
)

type DeviceGroupCommand struct {
	Command     string            `json:"command"`
	Id          string            `json:"id"`
	Owner       string            `json:"owner"`
	DeviceGroup model.DeviceGroup `json:"device_group"`
}

func (this *Publisher) PublishDeviceGroup(group model.DeviceGroup, userId string) (err error) {
	cmd := DeviceGroupCommand{Command: "PUT", Id: group.Id, DeviceGroup: group, Owner: userId}
	return this.PublishDeviceGroupCommand(cmd)
}

func (this *Publisher) PublishDeviceGroupDelete(id string, userId string) error {
	cmd := DeviceGroupCommand{Command: "DELETE", Id: id, Owner: userId}
	return this.PublishDeviceGroupCommand(cmd)
}

func (this *Publisher) PublishDeviceGroupCommand(cmd DeviceGroupCommand) error {
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce deviceGroup", cmd)
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.devicegroups.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(cmd.Id),
			Value: message,
			Time:  time.Now(),
		},
	)
	if err != nil {
		debug.PrintStack()
	}
	return err
}
