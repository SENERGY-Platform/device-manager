/*
 * Copyright waitDoneTrys24 InfAI (CC SES)
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
	"fmt"
	"github.com/SENERGY-Platform/device-manager/lib/api"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller"
	"github.com/SENERGY-Platform/device-manager/lib/tests/docker"
	"github.com/SENERGY-Platform/device-manager/lib/tests/helper"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

const waitDoneTrys = 20

func TestWaitDone(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.Load("./../../config.json")
	if err != nil {
		t.Fatal("ERROR: unable to load config", err)
	}

	port, err := helper.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	conf.ServerPort = strconv.Itoa(port)

	conf.DeviceRepoUrl, conf.PermissionsV2Url, conf.KafkaUrl, err = docker.DeviceRepoWithDependencies(ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}

	ctrl, err := controller.New(ctx, conf)
	if err != nil {
		t.Fatal(err)
	}

	srv, err := api.Start(conf, ctrl)
	if err != nil {
		t.Fatal("ERROR: unable to start api", err)
	}
	defer srv.Shutdown(context.Background())

	time.Sleep(2 * time.Second)

	signal.Sub("", signal.Known.UpdateDone, func(value string, wg *sync.WaitGroup) {
		log.Printf("TEST-DEBUG: received update done %#v\n", value)
	})

	t.Run("aspects", func(t *testing.T) {
		resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/aspects", models.Aspect{Id: a1Id, Name: a1Id})
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			t.Fatal(resp.Status, resp.StatusCode, string(b))
		}
		resp.Body.Close()
	})
	t.Run("functions", func(t *testing.T) {
		resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/functions", models.Function{Id: f1Id, Name: f1Id})
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			t.Fatal(resp.Status, resp.StatusCode, string(b))
		}
		resp.Body.Close()
	})

	protocol := models.Protocol{}
	t.Run("create protocol", func(t *testing.T) {
		resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/protocols", models.Protocol{
			Name:             "pname1",
			Handler:          "ph1",
			ProtocolSegments: []models.ProtocolSegment{{Name: "ps1"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			t.Fatal(resp.Status, resp.StatusCode, string(b))
		}

		err = json.NewDecoder(resp.Body).Decode(&protocol)
		if err != nil {
			t.Fatal(err)
		}
	})

	time.Sleep(2 * time.Second)

	t.Run("create device-types parallel", func(t *testing.T) {
		for i := range waitDoneTrys {
			t.Run(fmt.Sprintf("check device-type %v", i), func(t *testing.T) {
				t.Parallel()
				dt := models.DeviceType{}
				t.Run(fmt.Sprintf("create device-type %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/device-types?wait=true", models.DeviceType{
						Name:          fmt.Sprintf("foo-%v", i),
						DeviceClassId: "dc1",
						Services: []models.Service{
							{
								Name:    "s1name",
								LocalId: "lid1",
								Inputs: []models.Content{
									{
										ProtocolSegmentId: protocol.ProtocolSegments[0].Id,
										Serialization:     "json",
										ContentVariable: models.ContentVariable{
											Name:       "v1name",
											Type:       models.String,
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
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&dt)
					if err != nil {
						t.Fatal(err)
					}

					if dt.Id == "" {
						t.Fatal(dt)
					}
				})

				t.Run(fmt.Sprintf("read device-type %v", i), func(t *testing.T) {
					result := models.DeviceType{}
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/device-types/"+url.PathEscape(dt.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Log("http://localhost:" + conf.ServerPort + "/device-types/" + url.PathEscape(dt.Id))
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					result = models.DeviceType{}
					err = json.NewDecoder(resp.Body).Decode(&result)
					if err != nil {
						t.Fatal(err)
					}

					if result.Name != fmt.Sprintf("foo-%v", i) ||
						result.DeviceClassId != "dc1" ||
						len(result.Services) != 1 ||
						result.Services[0].Name != "s1name" ||
						result.Services[0].ProtocolId != protocol.Id ||
						result.Services[0].Inputs[0].ContentVariable.AspectId != a1Id ||
						result.Services[0].Inputs[0].ContentVariable.FunctionId != f1Id {
						t.Fatal(result)
					}
				})

				t.Run(fmt.Sprintf("delete device-type %v", i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/device-types/"+url.PathEscape(dt.Id)+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read device-type %v after delete", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/device-types/"+url.PathEscape(dt.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	})

	t.Run("create device-types", func(t *testing.T) {
		for i := range waitDoneTrys {
			i = i + waitDoneTrys
			t.Run(fmt.Sprintf("check device-type %v", i), func(t *testing.T) {
				dt := models.DeviceType{}
				t.Run(fmt.Sprintf("create device-type %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/device-types?wait=true", models.DeviceType{
						Name:          fmt.Sprintf("foo-%v", i),
						DeviceClassId: "dc1",
						Services: []models.Service{
							{
								Name:    "s1name",
								LocalId: "lid1",
								Inputs: []models.Content{
									{
										ProtocolSegmentId: protocol.ProtocolSegments[0].Id,
										Serialization:     "json",
										ContentVariable: models.ContentVariable{
											Name:       "v1name",
											Type:       models.String,
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
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&dt)
					if err != nil {
						t.Fatal(err)
					}

					if dt.Id == "" {
						t.Fatal(dt)
					}
				})

				t.Run(fmt.Sprintf("read device-type %v", i), func(t *testing.T) {
					result := models.DeviceType{}
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/device-types/"+url.PathEscape(dt.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Log("http://localhost:" + conf.ServerPort + "/device-types/" + url.PathEscape(dt.Id))
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					result = models.DeviceType{}
					err = json.NewDecoder(resp.Body).Decode(&result)
					if err != nil {
						t.Fatal(err)
					}

					if result.Name != fmt.Sprintf("foo-%v", i) ||
						result.DeviceClassId != "dc1" ||
						len(result.Services) != 1 ||
						result.Services[0].Name != "s1name" ||
						result.Services[0].ProtocolId != protocol.Id ||
						result.Services[0].Inputs[0].ContentVariable.AspectId != a1Id ||
						result.Services[0].Inputs[0].ContentVariable.FunctionId != f1Id {
						t.Fatal(result)
					}
				})

				t.Run(fmt.Sprintf("delete device-type %v", i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/device-types/"+url.PathEscape(dt.Id)+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read device-type %v after delete", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/device-types/"+url.PathEscape(dt.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	})

	dt := models.DeviceType{}
	t.Run(fmt.Sprintf("create device-type for device test"), func(t *testing.T) {
		resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/device-types?wait=true", models.DeviceType{
			Name:          "foo",
			DeviceClassId: "dc1",
			Services: []models.Service{
				{
					Name:    "s1name",
					LocalId: "lid1",
					Inputs: []models.Content{
						{
							ProtocolSegmentId: protocol.ProtocolSegments[0].Id,
							Serialization:     "json",
							ContentVariable: models.ContentVariable{
								Name:       "v1name",
								Type:       models.String,
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
			b, _ := io.ReadAll(resp.Body)
			t.Fatal(resp.Status, resp.StatusCode, string(b))
		}

		err = json.NewDecoder(resp.Body).Decode(&dt)
		if err != nil {
			t.Fatal(err)
		}

		if dt.Id == "" {
			t.Fatal(dt)
		}
	})

	t.Run("create devices parallel", func(t *testing.T) {
		for i := range waitDoneTrys {
			t.Run(fmt.Sprintf("check device %v", i), func(t *testing.T) {
				t.Parallel()
				device := models.DeviceType{}
				t.Run(fmt.Sprintf("create device %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/devices?wait=true", models.Device{
						LocalId:      fmt.Sprintf("foo-%v", i),
						Name:         fmt.Sprintf("foo-%v", i),
						DeviceTypeId: dt.Id,
					})
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&device)
					if err != nil {
						t.Fatal(err)
					}

					if device.Id == "" {
						t.Fatal(device)
					}
				})

				t.Run(fmt.Sprintf("read device %v", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/devices/"+url.PathEscape(device.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					result := models.Device{}
					err = json.NewDecoder(resp.Body).Decode(&result)
					if err != nil {
						t.Fatal(err)
					}

					if result.Name != fmt.Sprintf("foo-%v", i) {
						t.Fatal(result)
					}
				})

				t.Run(fmt.Sprintf("delete device-type %v", i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/devices/"+url.PathEscape(device.Id)+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read device %v after delete", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/devices/"+url.PathEscape(device.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	})

	t.Run("create devices", func(t *testing.T) {
		for i := range waitDoneTrys {
			i = i + waitDoneTrys
			t.Run(fmt.Sprintf("check device %v", i), func(t *testing.T) {
				device := models.DeviceType{}
				t.Run(fmt.Sprintf("create device %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/devices?wait=true", models.Device{
						LocalId:      fmt.Sprintf("foo-%v", i),
						Name:         fmt.Sprintf("foo-%v", i),
						DeviceTypeId: dt.Id,
					})
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&device)
					if err != nil {
						t.Fatal(err)
					}

					if device.Id == "" {
						t.Fatal(device)
					}
				})

				t.Run(fmt.Sprintf("read device %v", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/devices/"+url.PathEscape(device.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					result := models.Device{}
					err = json.NewDecoder(resp.Body).Decode(&result)
					if err != nil {
						t.Fatal(err)
					}

					if result.Name != fmt.Sprintf("foo-%v", i) {
						t.Fatal(result)
					}
				})

				t.Run(fmt.Sprintf("delete device-type %v", i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/devices/"+url.PathEscape(device.Id)+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read device %v after delete", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/devices/"+url.PathEscape(device.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	})

	t.Run("create hubs parallel", func(t *testing.T) {
		for i := range waitDoneTrys {
			t.Run(fmt.Sprintf("check hub %v", i), func(t *testing.T) {
				t.Parallel()
				hub := models.Hub{}
				t.Run(fmt.Sprintf("create hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/hubs?wait=true", models.Hub{
						Name: fmt.Sprintf("foo-%v", i),
						Hash: fmt.Sprintf("foo-%v", i),
					})
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&hub)
					if err != nil {
						t.Fatal(err)
					}

					if hub.Id == "" {
						t.Fatal(hub)
					}
				})

				t.Run(fmt.Sprintf("read hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/hubs/"+url.PathEscape(hub.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					result := models.Hub{}
					err = json.NewDecoder(resp.Body).Decode(&result)
					if err != nil {
						t.Fatal(err)
					}

					if result.Name != fmt.Sprintf("foo-%v", i) {
						t.Fatal(result)
					}
					if result.Hash != fmt.Sprintf("foo-%v", i) {
						t.Fatal(result)
					}
				})

				t.Run(fmt.Sprintf("delete hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/hubs/"+url.PathEscape(hub.Id)+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read hub %v after delete", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/hubs/"+url.PathEscape(hub.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	})

	t.Run("create hubs", func(t *testing.T) {
		for i := range waitDoneTrys {
			i = i + waitDoneTrys
			t.Run(fmt.Sprintf("check hub %v", i), func(t *testing.T) {
				hub := models.Hub{}
				t.Run(fmt.Sprintf("create hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/hubs?wait=true", models.Hub{
						Name: fmt.Sprintf("foo-%v", i),
						Hash: fmt.Sprintf("foo-%v", i),
					})
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&hub)
					if err != nil {
						t.Fatal(err)
					}

					if hub.Id == "" {
						t.Fatal(hub)
					}
				})

				t.Run(fmt.Sprintf("read hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/hubs/"+url.PathEscape(hub.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					result := models.Hub{}
					err = json.NewDecoder(resp.Body).Decode(&result)
					if err != nil {
						t.Fatal(err)
					}

					if result.Name != fmt.Sprintf("foo-%v", i) {
						t.Fatal(result)
					}
					if result.Hash != fmt.Sprintf("foo-%v", i) {
						t.Fatal(result)
					}
				})

				t.Run(fmt.Sprintf("delete hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/hubs/"+url.PathEscape(hub.Id)+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read hub %v after delete", i), func(t *testing.T) {
					resp, err := helper.Jwtget(userjwt, "http://localhost:"+conf.ServerPort+"/hubs/"+url.PathEscape(hub.Id))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	})

	t.Run("create local-devices parallel", testWait(conf, "local-devices", true, userjwt, func(i int) models.Device {
		return models.Device{
			Name:         fmt.Sprintf("l-foo-%v", i),
			LocalId:      fmt.Sprintf("l-foo-%v", i),
			DeviceTypeId: dt.Id,
		}
	}, func(e models.Device) string {
		return e.LocalId
	}, func(e models.Device, i int) error {
		if e.Name != fmt.Sprintf("l-foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("l-foo-%v", i))
		}
		return nil
	}))

	t.Run("create local-devices", testWait(conf, "local-devices", false, userjwt, func(i int) models.Device {
		return models.Device{
			Name:         fmt.Sprintf("l-foo-%v", i),
			LocalId:      fmt.Sprintf("l-foo-%v", i),
			DeviceTypeId: dt.Id,
		}
	}, func(e models.Device) string {
		return e.LocalId
	}, func(e models.Device, i int) error {
		if e.Name != fmt.Sprintf("l-foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("l-foo-%v", i))
		}
		return nil
	}))

	t.Run("create aspects", testWait(conf, "aspects", false, adminjwt, func(i int) models.Aspect {
		return models.Aspect{
			Name: fmt.Sprintf("foo-%v", i),
			SubAspects: []models.Aspect{
				{
					Name: fmt.Sprintf("foo-%v-0", i),
					SubAspects: []models.Aspect{
						{
							Name: fmt.Sprintf("foo-%v-0-1", i),
						},
						{
							Name: fmt.Sprintf("foo-%v-0-2", i),
						},
					},
				},
				{
					Name: fmt.Sprintf("foo-%v-1", i),
					SubAspects: []models.Aspect{
						{
							Name: fmt.Sprintf("foo-%v-1-1", i),
						},
						{
							Name: fmt.Sprintf("foo-%v-1-2", i),
						},
					},
				},
			},
		}
	}, func(aspect models.Aspect) string {
		return aspect.Id
	}, func(aspect models.Aspect, i int) error {
		if aspect.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", aspect.Name, fmt.Sprintf("foo-%v", i))
		}
		if len(aspect.SubAspects) != 2 {
			return fmt.Errorf("len SubAspects does not match expected")
		}
		return nil
	}))

	t.Run("create aspects parallel", testWait(conf, "aspects", true, adminjwt, func(i int) models.Aspect {
		return models.Aspect{
			Name: fmt.Sprintf("foo-%v", i),
			SubAspects: []models.Aspect{
				{
					Name: fmt.Sprintf("foo-%v-0", i),
					SubAspects: []models.Aspect{
						{
							Name: fmt.Sprintf("foo-%v-0-1", i),
						},
						{
							Name: fmt.Sprintf("foo-%v-0-2", i),
						},
					},
				},
				{
					Name: fmt.Sprintf("foo-%v-1", i),
					SubAspects: []models.Aspect{
						{
							Name: fmt.Sprintf("foo-%v-1-1", i),
						},
						{
							Name: fmt.Sprintf("foo-%v-1-2", i),
						},
					},
				},
			},
		}
	}, func(aspect models.Aspect) string {
		return aspect.Id
	}, func(aspect models.Aspect, i int) error {
		if aspect.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", aspect.Name, fmt.Sprintf("foo-%v", i))
		}
		if len(aspect.SubAspects) != 2 {
			return fmt.Errorf("len SubAspects does not match expected")
		}
		return nil
	}))

	t.Run("create functions parallel", testWait(conf, "functions", true, adminjwt, func(i int) models.Function {
		return models.Function{
			Id:          fmt.Sprintf("%vcontrolling-function:foo-%v", models.URN_PREFIX, i),
			Name:        fmt.Sprintf("foo-%v", i),
			Description: fmt.Sprintf("foo-%v", i),
		}
	}, func(e models.Function) string {
		return e.Id
	}, func(e models.Function, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Description != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("description %v does not match expected %v", e.Description, fmt.Sprintf("foo-%v", i))
		}
		return nil
	}))

	t.Run("create functions", testWait(conf, "functions", false, adminjwt, func(i int) models.Function {
		return models.Function{
			Id:          fmt.Sprintf("%vcontrolling-function:foo-%v", models.URN_PREFIX, i),
			Name:        fmt.Sprintf("foo-%v", i),
			Description: fmt.Sprintf("foo-%v", i),
		}
	}, func(e models.Function) string {
		return e.Id
	}, func(e models.Function, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Description != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("description %v does not match expected %v", e.Description, fmt.Sprintf("foo-%v", i))
		}
		return nil
	}))

	t.Run("create concepts parallel", testWait(conf, "concepts", true, adminjwt, func(i int) models.Concept {
		return models.Concept{
			Name: fmt.Sprintf("foo-%v", i),
		}
	}, func(e models.Concept) string {
		return e.Id
	}, func(e models.Concept, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		return nil
	}))

	t.Run("create concepts", testWait(conf, "concepts", false, adminjwt, func(i int) models.Concept {
		return models.Concept{
			Name: fmt.Sprintf("foo-%v", i),
		}
	}, func(e models.Concept) string {
		return e.Id
	}, func(e models.Concept, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		return nil
	}))

	t.Run("create characteristics parallel", testWait(conf, "characteristics", true, adminjwt, func(i int) models.Characteristic {
		return models.Characteristic{
			Name: fmt.Sprintf("foo-%v", i),
			Type: models.Structure,
			SubCharacteristics: []models.Characteristic{
				{
					Name: "foo",
					Type: models.String,
				},
				{
					Name: "bar",
					Type: models.String,
				},
			},
		}
	}, func(e models.Characteristic) string {
		return e.Id
	}, func(e models.Characteristic, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if len(e.SubCharacteristics) != 2 {
			return fmt.Errorf("len SubCharacteristics does not match expected")
		}
		return nil
	}))

	t.Run("create characteristics", testWait(conf, "characteristics", false, adminjwt, func(i int) models.Characteristic {
		return models.Characteristic{
			Name: fmt.Sprintf("foo-%v", i),
			Type: models.Structure,
			SubCharacteristics: []models.Characteristic{
				{
					Name: "foo",
					Type: models.String,
				},
				{
					Name: "bar",
					Type: models.String,
				},
			},
		}
	}, func(e models.Characteristic) string {
		return e.Id
	}, func(e models.Characteristic, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if len(e.SubCharacteristics) != 2 {
			return fmt.Errorf("len SubCharacteristics does not match expected")
		}
		return nil
	}))

	t.Run("create device-classes parallel", testWait(conf, "device-classes", true, adminjwt, func(i int) models.DeviceClass {
		return models.DeviceClass{
			Name:  fmt.Sprintf("foo-%v", i),
			Image: "http://foobar.foo/foo.jpg",
		}
	}, func(e models.DeviceClass) string {
		return e.Id
	}, func(e models.DeviceClass, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Image != "http://foobar.foo/foo.jpg" {
			return fmt.Errorf("image %v does not match expected %v", e.Image, "http://foobar.foo/foo.jpg")
		}
		return nil
	}))

	t.Run("create device-classes", testWait(conf, "device-classes", false, adminjwt, func(i int) models.DeviceClass {
		return models.DeviceClass{
			Name:  fmt.Sprintf("foo-%v", i),
			Image: "http://foobar.foo/foo.jpg",
		}
	}, func(e models.DeviceClass) string {
		return e.Id
	}, func(e models.DeviceClass, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Image != "http://foobar.foo/foo.jpg" {
			return fmt.Errorf("image %v does not match expected %v", e.Image, "http://foobar.foo/foo.jpg")
		}
		return nil
	}))

	t.Run("create device-groups parallel", testWait(conf, "device-groups", true, adminjwt, func(i int) models.DeviceGroup {
		return models.DeviceGroup{
			Name:  fmt.Sprintf("foo-%v", i),
			Image: "http://foobar.foo/foo.jpg",
		}
	}, func(e models.DeviceGroup) string {
		return e.Id
	}, func(e models.DeviceGroup, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Image != "http://foobar.foo/foo.jpg" {
			return fmt.Errorf("image %v does not match expected %v", e.Image, "http://foobar.foo/foo.jpg")
		}
		return nil
	}))

	t.Run("create device-groups", testWait(conf, "device-groups", false, adminjwt, func(i int) models.DeviceGroup {
		return models.DeviceGroup{
			Name:  fmt.Sprintf("foo-%v", i),
			Image: "http://foobar.foo/foo.jpg",
		}
	}, func(e models.DeviceGroup) string {
		return e.Id
	}, func(e models.DeviceGroup, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Image != "http://foobar.foo/foo.jpg" {
			return fmt.Errorf("image %v does not match expected %v", e.Image, "http://foobar.foo/foo.jpg")
		}
		return nil
	}))

	t.Run("create locations parallel", testWait(conf, "locations", true, adminjwt, func(i int) models.Location {
		return models.Location{
			Name:  fmt.Sprintf("foo-%v", i),
			Image: "http://foobar.foo/foo.jpg",
		}
	}, func(e models.Location) string {
		return e.Id
	}, func(e models.Location, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Image != "http://foobar.foo/foo.jpg" {
			return fmt.Errorf("image %v does not match expected %v", e.Image, "http://foobar.foo/foo.jpg")
		}
		return nil
	}))

	t.Run("create locations", testWait(conf, "locations", false, adminjwt, func(i int) models.Location {
		return models.Location{
			Name:  fmt.Sprintf("foo-%v", i),
			Image: "http://foobar.foo/foo.jpg",
		}
	}, func(e models.Location) string {
		return e.Id
	}, func(e models.Location, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		if e.Image != "http://foobar.foo/foo.jpg" {
			return fmt.Errorf("image %v does not match expected %v", e.Image, "http://foobar.foo/foo.jpg")
		}
		return nil
	}))

	t.Run("create protocols parallel", testWait(conf, "protocols", true, adminjwt, func(i int) models.Protocol {
		return models.Protocol{
			Name:    fmt.Sprintf("foo-%v", i),
			Handler: "ph1",
			ProtocolSegments: []models.ProtocolSegment{
				{Name: fmt.Sprintf("ps1-foo-%v", i)},
			},
		}
	}, func(e models.Protocol) string {
		return e.Id
	}, func(e models.Protocol, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		return nil
	}))

	t.Run("create protocols", testWait(conf, "protocols", false, adminjwt, func(i int) models.Protocol {
		return models.Protocol{
			Name:    fmt.Sprintf("foo-%v", i),
			Handler: "ph1",
			ProtocolSegments: []models.ProtocolSegment{
				{Name: fmt.Sprintf("ps1-foo-%v", i)},
			},
		}
	}, func(e models.Protocol) string {
		return e.Id
	}, func(e models.Protocol, i int) error {
		if e.Name != fmt.Sprintf("foo-%v", i) {
			return fmt.Errorf("name %v does not match expected %v", e.Name, fmt.Sprintf("foo-%v", i))
		}
		return nil
	}))
}

func testWait[T any](conf config.Config, resource string, parallel bool, token string, create func(i int) T, getId func(T) string, check func(T, int) error) func(t *testing.T) {
	return func(t *testing.T) {
		for i := range waitDoneTrys {
			if !parallel {
				i = i + waitDoneTrys
			}
			t.Run(fmt.Sprintf("check %v %v", resource, i), func(t *testing.T) {
				if parallel {
					t.Parallel()
				}
				element := create(i)
				t.Run(fmt.Sprintf("create %v %v", resource, i), func(t *testing.T) {
					resp, err := helper.Jwtpost(token, "http://localhost:"+conf.ServerPort+"/"+resource+"?wait=true", element)
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					err = json.NewDecoder(resp.Body).Decode(&element)
					if err != nil {
						t.Fatal(err)
					}

					err = check(element, i)
					if err != nil {
						t.Fatal(err)
					}
				})

				t.Run(fmt.Sprintf("read %v %v", resource, i), func(t *testing.T) {
					resp, err := helper.Jwtget(token, "http://localhost:"+conf.ServerPort+"/"+resource+"/"+url.PathEscape(getId(element)))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}

					var element T
					err = json.NewDecoder(resp.Body).Decode(&element)
					if err != nil {
						t.Fatal(err)
					}

					err = check(element, i)
					if err != nil {
						t.Fatal(err)
					}
				})

				t.Run(fmt.Sprintf("delete %v %v", resource, i), func(t *testing.T) {
					resp, err := helper.Jwtdelete(token, "http://localhost:"+conf.ServerPort+"/"+resource+"/"+url.PathEscape(getId(element))+"?wait=true")
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})

				t.Run(fmt.Sprintf("read %v %v after delete", resource, i), func(t *testing.T) {
					resp, err := helper.Jwtget(token, "http://localhost:"+conf.ServerPort+"/"+resource+"/"+url.PathEscape(getId(element)))
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						b, _ := io.ReadAll(resp.Body)
						t.Fatal(resp.Status, resp.StatusCode, string(b))
					}
				})
			})

		}
	}
}
