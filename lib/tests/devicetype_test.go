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

package tests

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/tests/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func testDeviceType(t *testing.T, port string) {
	resp, err := helper.Jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtpost(adminjwt, "http://localhost:"+port+"/protocols", model.Protocol{
		Name:             "pname1",
		Handler:          "ph1",
		ProtocolSegments: []model.ProtocolSegment{{Name: "ps1"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	protocol := model.Protocol{}
	err = json.NewDecoder(resp.Body).Decode(&protocol)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{
		Name:          "foo",
		DeviceClassId: "dc1",
		Services: []model.Service{
			{
				Name:    "s1name",
				LocalId: "lid1",
				Inputs: []model.Content{
					{
						ProtocolSegmentId: protocol.ProtocolSegments[0].Id,
						Serialization:     "json",
						ContentVariable: model.ContentVariable{
							Name:       "v1name",
							Type:       model.String,
							FunctionId: "f1",
							AspectId:   "a1",
						},
					},
				},
				ProtocolId: protocol.Id,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	dt := model.DeviceType{}
	err = json.NewDecoder(resp.Body).Decode(&dt)
	if err != nil {
		t.Fatal(err)
	}

	if dt.Id == "" {
		t.Fatal(dt)
	}

	result := model.DeviceType{}
	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Log("http://localhost:" + port + "/device-types/" + url.PathEscape(dt.Id))
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	result = model.DeviceType{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "foo" ||
		result.DeviceClassId != "dc1" ||
		len(result.Services) != 1 ||
		result.Services[0].Name != "s1name" ||
		result.Services[0].ProtocolId != protocol.Id ||
		result.Services[0].Inputs[0].ContentVariable.AspectId != "a1" ||
		result.Services[0].Inputs[0].ContentVariable.FunctionId != "f1" {
		t.Fatal(result)
	}

	resp, err = helper.Jwtdelete(adminjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}
}

func testDeviceTypeWithServiceGroups(t *testing.T, port string) {
	resp, err := helper.Jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtpost(adminjwt, "http://localhost:"+port+"/protocols", model.Protocol{
		Name:             "pname2",
		Handler:          "ph2",
		ProtocolSegments: []model.ProtocolSegment{{Name: "ps2"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	protocol := model.Protocol{}
	err = json.NewDecoder(resp.Body).Decode(&protocol)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{
		Name:          "foo",
		DeviceClassId: "dc1",
		ServiceGroups: []model.ServiceGroup{
			{
				Key:         "sg1",
				Name:        "service group 1",
				Description: "foo  bar",
			},
		},
		Services: []model.Service{
			{
				Name:    "s1name",
				LocalId: "lid1",
				Inputs: []model.Content{
					{
						ProtocolSegmentId: protocol.ProtocolSegments[0].Id,
						Serialization:     "json",
						ContentVariable: model.ContentVariable{
							Name:       "v1name",
							Type:       model.String,
							FunctionId: "f1",
							AspectId:   "a1",
						},
					},
				},
				ProtocolId: protocol.Id,
			},
			{
				Name:    "s2name",
				LocalId: "lid2",
				Inputs: []model.Content{
					{
						ProtocolSegmentId: protocol.ProtocolSegments[0].Id,
						Serialization:     "json",
						ContentVariable: model.ContentVariable{
							Name:       "v1name",
							Type:       model.String,
							FunctionId: "f1",
							AspectId:   "a1",
						},
					},
				},
				ProtocolId:      protocol.Id,
				ServiceGroupKey: "sg1",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	dt := model.DeviceType{}
	err = json.NewDecoder(resp.Body).Decode(&dt)
	if err != nil {
		t.Fatal(err)
	}

	if dt.Id == "" {
		t.Fatal(dt)
	}

	result := model.DeviceType{}
	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Log("http://localhost:" + port + "/device-types/" + url.PathEscape(dt.Id))
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	result = model.DeviceType{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "foo" ||
		result.DeviceClassId != "dc1" ||
		len(result.Services) != 2 {
		t.Fatal(result.Name, result.DeviceClassId, len(result.Services))
	}

	if !reflect.DeepEqual(result.ServiceGroups, []model.ServiceGroup{
		{
			Key:         "sg1",
			Name:        "service group 1",
			Description: "foo  bar",
		},
	}) {
		t.Fatal(result.ServiceGroups)
	}

	if result.Services[0].Name != "s1name" ||
		result.Services[0].LocalId != "lid1" ||
		result.Services[0].ServiceGroupKey != "" ||
		result.Services[0].ProtocolId != protocol.Id ||
		result.Services[0].Inputs[0].ContentVariable.AspectId != "a1" ||
		result.Services[0].Inputs[0].ContentVariable.FunctionId != "f1" {

		t.Fatal(result.Services[0])
	}

	if result.Services[1].Name != "s2name" ||
		result.Services[1].LocalId != "lid2" ||
		result.Services[1].ServiceGroupKey != "sg1" ||
		result.Services[1].ProtocolId != protocol.Id ||
		result.Services[1].Inputs[0].ContentVariable.AspectId != "a1" ||
		result.Services[1].Inputs[0].ContentVariable.FunctionId != "f1" {
		temp, _ := json.Marshal(result.Services[1])
		t.Fatal(string(temp))
	}

	resp, err = helper.Jwtdelete(adminjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}
}
