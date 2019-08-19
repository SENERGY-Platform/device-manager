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

type ConceptCommand struct {
	Command    string           `json:"command"`
	Id         string           `json:"id"`
	Owner      string           `json:"owner"`
	Concept    model.Concept    `json:"concept"`
}

func (this *Publisher) PublishConcept(concept model.Concept, userId string) (err error) {
	cmd := ConceptCommand{Command: "PUT", Id: concept.Id, Concept: concept, Owner: userId}
	return this.PublishConceptCommand(cmd)
}

func (this *Publisher) PublishConceptDelete(id string, userId string) error {
	cmd := ConceptCommand{Command: "DELETE", Id: id, Owner: userId}
	return this.PublishConceptCommand(cmd)
}


func (this *Publisher) PublishConceptCommand(cmd ConceptCommand) error {
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce devicetype", cmd)
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = this.devicetypes.WriteMessages(
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
