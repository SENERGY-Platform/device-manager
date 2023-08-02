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
	"errors"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/segmentio/kafka-go"
	"log"
	"runtime/debug"
	"time"
)

type FunctionCommand struct {
	Command  string          `json:"command"`
	Id       string          `json:"id"`
	Owner    string          `json:"owner"`
	Function models.Function `json:"function"`
}

func (this *Publisher) PublishFunction(function models.Function, userId string) (err error) {
	cmd := FunctionCommand{Command: "PUT", Id: function.Id, Function: function, Owner: userId}
	return this.PublishFunctionCommand(cmd)
}

func (this *Publisher) PublishFunctionDelete(id string, userId string) error {
	cmd := FunctionCommand{Command: "DELETE", Id: id, Owner: userId}
	return this.PublishFunctionCommand(cmd)
}

func (this *Publisher) PublishFunctionCommand(cmd FunctionCommand) error {
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce function", cmd)
	}
	if cmd.Owner == "" {
		return errors.New("missing owner in command")
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.functions.WriteMessages(
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
