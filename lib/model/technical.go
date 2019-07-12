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

package model

import "github.com/google/uuid"

type Hub struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Hash           string   `json:"hash"`
	DeviceLocalIds []string `json:"device_local_ids"`
}

func (hub *Hub) GenerateId() {
	hub.Id = "urn:infai:ses:hub:" + uuid.New().String()
}

type Protocol struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Handler          string            `json:"handler"`
	ProtocolSegments []ProtocolSegment `json:"protocol_segments"`
}

func (protocol *Protocol) GenerateId() {
	protocol.Id = "urn:infai:ses:protocol:" + uuid.New().String()
	for i, segment := range protocol.ProtocolSegments {
		segment.GenerateId()
		protocol.ProtocolSegments[i] = segment
	}
}

type ProtocolSegment struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (segment *ProtocolSegment) GenerateId() {
	segment.Id = "urn:infai:ses:segment:" + uuid.New().String()
}

type Content struct {
	Id                   string                `json:"id"`
	Variable             Variable              `json:"variable"`
	Serialization        string                `json:"serialization"`
	SerializationOptions []SerializationOption `json:"serialization_options"`
	ProtocolSegmentId    string                `json:"protocol_segment_id"`
}

func (content *Content) GenerateId() {
	content.Id = "urn:infai:ses:content:" + uuid.New().String()
	for i, option := range content.SerializationOptions {
		option.GenerateId()
		content.SerializationOptions[i] = option
	}
	content.Variable.GenerateId()
}

type SerializationOption struct {
	Id         string `json:"id"`
	Option     string `json:"option"`
	VariableId string `json:"variable_id"`
}

func (option *SerializationOption) GenerateId() {
	option.Id = "urn:infai:ses:option:" + uuid.New().String()
}
