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

type AspectCommand struct {
	Command string       `json:"command"`
	Id      string       `json:"id"`
	Owner   string       `json:"owner"`
	Aspect  model.Aspect `json:"aspect"`
}

func (this *Publisher) PublishAspect(aspect model.Aspect, userId string) (err error) {
	cmd := AspectCommand{Command: "PUT", Id: aspect.Id, Aspect: aspect, Owner: userId}
	return this.PublishAspectCommand(cmd)
}

func (this *Publisher) PublishAspectDelete(id string, userId string) error {
	cmd := AspectCommand{Command: "DELETE", Id: id, Owner: userId}
	return this.PublishAspectCommand(cmd)
}

func (this *Publisher) PublishAspectCommand(cmd AspectCommand) error {
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce aspect", cmd)
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.aspects.WriteMessages(
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
