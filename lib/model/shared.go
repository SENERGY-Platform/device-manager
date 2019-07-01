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
	device.Id = uuid.New().String()
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
	deviceType.Id = uuid.New().String()
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

type VariableType string

const (
	String  VariableType = "http://www.w3.org/2001/XMLSchema#string"
	Integer VariableType = "http://www.w3.org/2001/XMLSchema#integer"
	Float   VariableType = "http://www.w3.org/2001/XMLSchema#decimal"
	Boolean VariableType = "http://www.w3.org/2001/XMLSchema#boolean"

	Array     VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#Array"     //array with predefined length where each element can be of a different type
	Structure VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#structure" //object with predefined fields where each field can be of a different type
	Map       VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#map"       //object/map where each element has to be of the same type but the key can change
	List      VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#list"      //array where each element has to be of the same type but the length can change
)

type Variable struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Type         VariableType `json:"type"`
	SubVariables []Variable   `json:"sub_variables"`
}
