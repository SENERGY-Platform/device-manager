/*
 * Copyright 2024 InfAI (CC SES)
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

	conf.DeviceRepoUrl, conf.PermissionsUrl, conf.KafkaUrl, err = docker.DeviceRepoWithDependencies(ctx, wg)
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
		for i := range 20 {
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
		for i := range 20 {
			i = i + 20
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
		for i := range 20 {
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
		for i := range 20 {
			i = i + 20
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
		for i := range 20 {
			t.Run(fmt.Sprintf("check hub %v", i), func(t *testing.T) {
				t.Parallel()
				hub := models.HubEdit{}
				t.Run(fmt.Sprintf("create hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/hubs?wait=true", models.HubEdit{
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
		for i := range 20 {
			i = i + 20
			t.Run(fmt.Sprintf("check hub %v", i), func(t *testing.T) {
				hub := models.HubEdit{}
				t.Run(fmt.Sprintf("create hub %v", i), func(t *testing.T) {
					resp, err := helper.Jwtpost(userjwt, "http://localhost:"+conf.ServerPort+"/hubs?wait=true", models.HubEdit{
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

}
