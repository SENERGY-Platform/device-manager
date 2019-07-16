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
	"time"
)

func testHub(t *testing.T, port string) {
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

	time.Sleep(10 * time.Second)

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{
		Name: "foo",
		DeviceClass: model.DeviceClass{
			Id: "dc1",
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
							Name: "v1name",
							Type: "string",
						},
					},
				},
				Functions: []model.Function{
					{Id: "f1"},
				},
				Aspects: []model.Aspect{
					{Id: "a1"},
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

	time.Sleep(10 * time.Second)

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d2",
		DeviceTypeId: dt.Id,
		LocalId:      "lid2",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	device := model.Device{}
	err = json.NewDecoder(resp.Body).Decode(&device)
	if err != nil {
		t.Fatal(err)
	}

	if device.Id == "" {
		t.Fatal(device)
	}

	time.Sleep(10 * time.Second)

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/hubs", model.Hub{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/hubs", model.Hub{
		Name:           "h1",
		DeviceLocalIds: []string{"unknown"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/hubs", model.Hub{
		Name:           "h1",
		Hash:           "foobar",
		DeviceLocalIds: []string{device.LocalId},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	hub := model.Hub{}
	err = json.NewDecoder(resp.Body).Decode(&hub)
	if err != nil {
		t.Fatal(err)
	}

	if hub.Id == "" {
		t.Fatal(hub)
	}

	time.Sleep(15 * time.Second)

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(hub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	result := model.Hub{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "h1" || result.Hash != "foobar" || !reflect.DeepEqual(result.DeviceLocalIds, []string{device.LocalId}) {
		t.Fatal(result)
	}

	resp, err = helper.Jwtdelete(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(hub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	time.Sleep(5 * time.Second)

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(hub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect 404 error
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal(resp.Status, resp.StatusCode)
	}
}

func testHubAssertions(t *testing.T, port string) {
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

	time.Sleep(5 * time.Second)

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{
		Name: "foo",
		DeviceClass: model.DeviceClass{
			Id: "dc1",
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
							Name: "v1name",
							Type: "string",
						},
					},
				},
				Functions: []model.Function{
					{Id: "f1"},
				},
				Aspects: []model.Aspect{
					{Id: "a1"},
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

	time.Sleep(15 * time.Second)

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d3",
		DeviceTypeId: dt.Id,
		LocalId:      "lid3",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	d3 := model.Device{}
	err = json.NewDecoder(resp.Body).Decode(&d3)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d4",
		DeviceTypeId: dt.Id,
		LocalId:      "lid4",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	d4 := model.Device{}
	err = json.NewDecoder(resp.Body).Decode(&d4)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/devices", model.Device{
		Name:         "d5",
		DeviceTypeId: dt.Id,
		LocalId:      "lid5",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	time.Sleep(15 * time.Second)

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/hubs", model.Hub{
		Name:           "h2",
		Hash:           "foobar",
		DeviceLocalIds: []string{"lid3", "lid4", "lid5"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	hub := model.Hub{}
	err = json.NewDecoder(resp.Body).Decode(&hub)
	if err != nil {
		t.Fatal(err)
	}

	if hub.Id == "" {
		t.Fatal(hub)
	}

	time.Sleep(15 * time.Second)

	// update hub on device local id change

	resp, err = helper.Jwtput(userjwt, "http://localhost:"+port+"/devices/"+url.PathEscape(d3.Id), model.Device{
		Id:           d3.Id,
		Name:         "d3",
		DeviceTypeId: dt.Id,
		LocalId:      "lid3_changed",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	time.Sleep(15 * time.Second)
	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(hub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	result := model.Hub{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != hub.Name || result.Hash != "" || !reflect.DeepEqual(result.DeviceLocalIds, []string{"lid4", "lid5"}) {
		t.Fatal(result)
	}

	// update hub on device delete

	resp, err = helper.Jwtdelete(userjwt, "http://localhost:"+port+"/devices/"+url.PathEscape(d4.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	time.Sleep(15 * time.Second)
	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(hub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	result = model.Hub{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != hub.Name || result.Hash != "" || !reflect.DeepEqual(result.DeviceLocalIds, []string{"lid5"}) {
		t.Fatal(result)
	}

	// only one hub may have device

	resp, err = helper.Jwtpost(userjwt, "http://localhost:"+port+"/hubs", model.Hub{
		Name:           "h3",
		Hash:           "foobar",
		DeviceLocalIds: []string{"lid5"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	newHub := model.Hub{}
	err = json.NewDecoder(resp.Body).Decode(&newHub)
	if err != nil {
		t.Fatal(err)
	}

	if newHub.Id == "" {
		t.Fatal(newHub)
	}

	time.Sleep(20 * time.Second)

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(hub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	err = json.NewDecoder(resp.Body).Decode(&hub)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = helper.Jwtget(userjwt, "http://localhost:"+port+"/hubs/"+url.PathEscape(newHub.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	err = json.NewDecoder(resp.Body).Decode(&newHub)
	if err != nil {
		t.Fatal(err)
	}

	if len(hub.DeviceLocalIds) != 0 || len(newHub.DeviceLocalIds) != 1 {
		t.Fatal(hub, newHub)
	}
}
