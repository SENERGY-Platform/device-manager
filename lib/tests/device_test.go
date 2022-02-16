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
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/api"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/tests/helper"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"testing"
)

func testDevice(t *testing.T, port string) {
	resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+port+"/protocols", model.Protocol{
		Name:             "p2",
		Handler:          "ph1",
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
							FunctionId: f1Id,
							AspectId:   a1Id,
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

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:    "d1",
		LocalId: "lid1",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d1",
		DeviceTypeId: dt.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d1",
		DeviceTypeId: dt.Id,
		LocalId:      "lid1",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	device := model.Device{}
	err = json.NewDecoder(resp.Body).Decode(&device)
	if err != nil {
		t.Fatal(err)
	}

	if device.Id == "" {
		t.Fatal(device)
	}

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/devices/"+url.PathEscape(device.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	result := model.Device{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "d1" || result.LocalId != "lid1" || result.DeviceTypeId != dt.Id {
		t.Fatal(result)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "reused_local_id",
		DeviceTypeId: dt.Id,
		LocalId:      "lid1",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal("device.local_id should be validated for global uniqueness: ", resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtdelete(userjwt, "http://localhost:"+port+"/devices/"+url.PathEscape(device.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/devices/"+url.PathEscape(device.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect 404 error
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal(resp.Status, resp.StatusCode)
	}
}

func testDeviceAttributes(t *testing.T, port string) {
	resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+port+"/protocols", model.Protocol{
		Name:             "p2",
		Handler:          "ph1",
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
							FunctionId: f1Id,
							AspectId:   a1Id,
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

	deviceId, err := initDevice(port, dt)

	t.Run("normal attr init", tryDeviceAttributeUpdate(port, dt.Id, deviceId, "", []model.Attribute{
		{
			Key:    "a1",
			Value:  "va1",
			Origin: "",
		},
		{
			Key:    "a2",
			Value:  "va2",
			Origin: "test1",
		},
		{
			Key:    "a3",
			Value:  "va3",
			Origin: "test1",
		},
		{
			Key:    "a4",
			Value:  "va4",
			Origin: "test2",
		},
		{
			Key:    "a5",
			Value:  "va5",
			Origin: "test2",
		},
	}, []model.Attribute{
		{
			Key:    "a1",
			Value:  "va1",
			Origin: "",
		},
		{
			Key:    "a2",
			Value:  "va2",
			Origin: "test1",
		},
		{
			Key:    "a3",
			Value:  "va3",
			Origin: "test1",
		},
		{
			Key:    "a4",
			Value:  "va4",
			Origin: "test2",
		},
		{
			Key:    "a5",
			Value:  "va5",
			Origin: "test2",
		},
	}))

	t.Run("normal attr update", tryDeviceAttributeUpdate(port, dt.Id, deviceId, "", []model.Attribute{
		{
			Key:    "a12",
			Value:  "va12",
			Origin: "",
		},
		{
			Key:    "a22",
			Value:  "va22",
			Origin: "test1",
		},
		{
			Key:    "a32",
			Value:  "va32",
			Origin: "test1",
		},
		{
			Key:    "a42",
			Value:  "va42",
			Origin: "test2",
		},
		{
			Key:    "a52",
			Value:  "va52",
			Origin: "test2",
		},
	}, []model.Attribute{
		{
			Key:    "a12",
			Value:  "va12",
			Origin: "",
		},
		{
			Key:    "a22",
			Value:  "va22",
			Origin: "test1",
		},
		{
			Key:    "a32",
			Value:  "va32",
			Origin: "test1",
		},
		{
			Key:    "a42",
			Value:  "va42",
			Origin: "test2",
		},
		{
			Key:    "a52",
			Value:  "va52",
			Origin: "test2",
		},
	}))

	t.Run("origin attr update", tryDeviceAttributeUpdate(port, dt.Id, deviceId, "test1", []model.Attribute{
		{
			Key:    "a13",
			Value:  "va13",
			Origin: "",
		},
		{
			Key:    "a23",
			Value:  "va23",
			Origin: "test1",
		},
		{
			Key:    "a33",
			Value:  "va33",
			Origin: "test1",
		},
		{
			Key:    "a43",
			Value:  "va43",
			Origin: "test2",
		},
		{
			Key:    "a53",
			Value:  "va53",
			Origin: "test2",
		},
	}, []model.Attribute{
		{
			Key:    "a12",
			Value:  "va12",
			Origin: "",
		},
		{
			Key:    "a23",
			Value:  "va23",
			Origin: "test1",
		},
		{
			Key:    "a33",
			Value:  "va33",
			Origin: "test1",
		},
		{
			Key:    "a42",
			Value:  "va42",
			Origin: "test2",
		},
		{
			Key:    "a52",
			Value:  "va52",
			Origin: "test2",
		},
	}))
}

func tryDeviceAttributeUpdate(port string, dtId string, deviceId string, origin string, attributes []model.Attribute, expected []model.Attribute) func(t *testing.T) {
	return func(t *testing.T) {
		sort.Slice(attributes, func(i, j int) bool {
			return attributes[i].Key < attributes[j].Key
		})

		endpoint := "http://localhost:" + port + "/devices/" + url.PathEscape(deviceId)
		if origin != "" {
			endpoint = endpoint + "?" + url.Values{api.UpdateOnlySameOriginAttributesKey: {origin}}.Encode()
		}
		resp, err := helper.Jwtput(userjwt, endpoint, model.Device{
			Id:           deviceId,
			Name:         "d1",
			LocalId:      uuid.New().String(),
			DeviceTypeId: dtId,
			Attributes:   attributes,
		})
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(resp.Body)
			t.Fatal(resp.Status, resp.StatusCode, string(b))
		}

		device := model.Device{}
		err = json.NewDecoder(resp.Body).Decode(&device)
		if err != nil {
			t.Fatal(err)
		}

		if device.Id == "" {
			t.Fatal(device)
		}
		if !reflect.DeepEqual(device.Attributes, expected) {
			t.Error(device, expected)
			return
		}

		resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/devices/"+url.PathEscape(device.Id))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(resp.Body)
			t.Fatal(resp.Status, resp.StatusCode, string(b))
		}

		result := model.Device{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(result.Attributes, expected) {
			t.Error(device, expected)
			return
		}
	}
}

func initDevice(port string, dt model.DeviceType) (string, error) {
	resp, err := helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d1",
		LocalId:      uuid.New().String(),
		DeviceTypeId: dt.Id,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	result := model.Device{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Id, err
}
