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
	"github.com/segmentio/kafka-go"
	"log"
	"runtime/debug"
	"time"
)

type PermCommandMsg struct {
	Command  string `json:"command"`
	Kind     string
	Resource string
	User     string
	Group    string
	Right    string
}

func (this *Publisher) PublishDeleteUserRights(resource string, id string, userId string) error {
	cmd := PermCommandMsg{
		Command:  "DELETE",
		Kind:     resource,
		Resource: id,
		User:     userId,
	}
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce Location", cmd)
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.permissions.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(userId + "_" + resource + "_" + id),
			Value: message,
			Time:  time.Now(),
		},
	)
	if err != nil {
		debug.PrintStack()
	}
	return err
}
