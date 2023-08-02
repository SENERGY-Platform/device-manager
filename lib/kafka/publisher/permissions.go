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
	"context"
	"encoding/json"
	"errors"
	permmodel "github.com/SENERGY-Platform/permission-search/lib/model"
	"github.com/segmentio/kafka-go"
	"log"
	"runtime/debug"
	"time"
)

type CommandWithRights = permmodel.CommandWithRights

func (this *Publisher) PublishRights(kind string, id string, element permmodel.ResourceRightsBase) error {
	cmd := CommandWithRights{
		Command: "RIGHTS",
		Id:      id,
		Rights:  &element,
	}
	key := id + "/rights"
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	if this.config.LogLevel == "DEBUG" || this.config.Debug {
		log.Printf("DEBUG: produce rights: topic=%v, key=%v message=%v", kind, key, string(message))
	}
	var writer *kafka.Writer
	switch kind {
	case this.config.DeviceTopic:
		writer = this.devices
	case this.config.DeviceGroupTopic:
		writer = this.devicegroups
	case this.config.HubTopic:
		writer = this.hubs
	case this.config.LocationTopic:
		writer = this.locations
	default:
		debug.PrintStack()
		return errors.New("unknown kind for PublishDeleteUserRights()")
	}
	err = writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: message,
			Time:  time.Now(),
		},
	)
	if err != nil {
		debug.PrintStack()
	}
	return err
}
