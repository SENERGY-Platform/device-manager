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
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/api"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/tests/mock"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestDeviceTypeWithMock(t *testing.T) {
	conf, err := config.Load("./../../config.json")
	if err != nil {
		t.Fatal("ERROR: unable to load config", err)
	}
	publ, conf, close := mock.New(conf)
	defer close()

	port, err := getFreePort()
	if err != nil {
		t.Fatal(err)
	}
	conf.ServerPort = strconv.Itoa(port)

	ctrl, err := controller.NewWithPublisher(conf, publ)
	if err != nil {
		t.Fatal(err)
	}

	srv, err := api.Start(conf, ctrl)
	if err != nil {
		log.Fatal("ERROR: unable to start api", err)
	}
	defer srv.Shutdown(context.Background())

	time.Sleep(200 * time.Millisecond)

	t.Run("testDeviceType", func(t *testing.T) {
		testDeviceType(t, conf.ServerPort)
	})

}

func testDeviceType(t *testing.T, port string) {
	resp, err := jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	//expect validation error
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = jwtpost(userjwt, "http://localhost:"+port+"/device-types", model.DeviceType{
		Name: "foo",
		DeviceClass: model.DeviceClass{
			Id: "dc1",
		},
		Services: []model.Service{
			{
				Name: "s1name",
				Functions: []model.Function{
					{Id: "f1"},
				},
				Aspects: []model.Aspect{
					{Id: "a1"},
				},
				ProtocolId: "p1",
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
	err = userjwt.GetJSON("http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id), &result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "foo" ||
		result.DeviceClass.Id != "dc1" ||
		len(result.Services) != 1 ||
		result.Services[0].Name != "s1name" ||
		result.Services[0].ProtocolId != "p1" ||
		len(result.Services[0].Aspects) != 1 ||
		result.Services[0].Aspects[0].Id != "a1" ||
		len(result.Services[0].Functions) != 1 ||
		result.Services[0].Functions[0].Id != "f1" {

		t.Fatal(result)
	}

	resp, err = jwtdelete(userjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}

	resp, err = jwtget(userjwt, "http://localhost:"+port+"/device-types/"+url.PathEscape(dt.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatal(resp.Status, resp.StatusCode)
	}
}
