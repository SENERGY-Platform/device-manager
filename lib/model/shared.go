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

type Device struct {
	Id           string `json:"id"`
	LocalId      string `json:"local_id"`
	Name         string `json:"name"`
	DeviceTypeId string `json:"device_type_id"`
}

func (device *Device) GenerateId() {
	device.Id = "urn:infai:ses:device:" + uuid.New().String()
}

type DeviceType struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Services    []Service   `json:"services"`
	DeviceClass DeviceClass `json:"device_class"`
}

func (deviceType *DeviceType) GenerateId() {
	deviceType.Id = "urn:infai:ses:device-type:" + uuid.New().String()
	for i, service := range deviceType.Services {
		service.GenerateId()
		deviceType.Services[i] = service
	}
	if deviceType.DeviceClass.Id == "" {
		deviceType.DeviceClass.GenerateId()
	}
}

type Service struct {
	Id          string     `json:"id"`
	LocalId     string     `json:"local_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Aspects     []Aspect   `json:"aspects"`
	ProtocolId  string     `json:"protocol_id"`
	Inputs      []Content  `json:"inputs"`
	Outputs     []Content  `json:"outputs"`
	Functions   []Function `json:"functions"`
}

func (service *Service) GenerateId() {
	service.Id = "urn:infai:ses:service:" + uuid.New().String()
	for i, function := range service.Functions {
		if function.Id == "" {
			function.GenerateId()
			service.Functions[i] = function
		}
	}
	for i, aspect := range service.Aspects {
		if aspect.Id == "" {
			aspect.GenerateId()
			service.Aspects[i] = aspect
		}
	}
	for i, content := range service.Inputs {
		content.GenerateId()
		service.Inputs[i] = content
	}
	for i, content := range service.Outputs {
		content.GenerateId()
		service.Outputs[i] = content
	}
}

type VariableType string

const (
	String  VariableType = "http://www.w3.org/2001/XMLSchema#string"
	Integer VariableType = "http://www.w3.org/2001/XMLSchema#integer"
	Float   VariableType = "http://www.w3.org/2001/XMLSchema#decimal"
	Boolean VariableType = "http://www.w3.org/2001/XMLSchema#boolean"

	Collection VariableType = "http://www.w3.org/1999/02/22-rdf-syntax-ns#List"
)

type Variable struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Type         VariableType `json:"type"`
	SubVariables []Variable   `json:"sub_variables"`
	Property     Property     `json:"property"`
}

func (variable *Variable) GenerateId() {
	variable.Id = "urn:infai:ses:variable:" + uuid.New().String()
	for i, v := range variable.SubVariables {
		v.GenerateId()
		variable.SubVariables[i] = v
	}
	variable.Property.GenerateId()
}

type Property struct {
	Id       string      `json:"id"`
	Unit     string      `json:"unit"`
	Value    interface{} `json:"value"`
	MinValue float64     `json:"min_value"`
	MaxValue float64     `json:"max_value"`
}

func (property *Property) GenerateId() {
	property.Id = "urn:infai:ses:property:" + uuid.New().String()
}
