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
	"github.com/SENERGY-Platform/permission-search/lib/model"
	"github.com/SENERGY-Platform/service-commons/pkg/donewait"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRightsReplay(t *testing.T) {
	if testing.Short() {
		t.Skip("disabled in short mode")
	}
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.Load("./../../config.json")
	if err != nil {
		t.Fatal("ERROR: unable to load config", err)
	}
	conf.HandleDoneWait = true

	port, err := helper.GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}
	conf.ServerPort = strconv.Itoa(port)

	_, zkIp, err := docker.Zookeeper(ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	zookeeperUrl := zkIp + ":2181"

	conf.KafkaUrl, err = docker.Kafka(ctx, wg, zookeeperUrl)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(1 * time.Second)

	_, elasticIp, err := docker.OpenSearch(ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}

	permCmd, _, permIp, err := docker.PermissionSearchWithCmdCallback(ctx, wg, conf.KafkaUrl, elasticIp)
	if err != nil {
		t.Error(err)
		return
	}
	conf.PermissionsUrl = "http://" + permIp + ":8080"

	_, mongoIp, err := docker.MongoDB(ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	mongoUrl := "mongodb://" + mongoIp + ":27017"

	time.Sleep(1 * time.Second)

	initialDeviceRepoCtx, initialDeviceRepoCancel := context.WithCancel(ctx)

	_, repoIp, err := docker.DeviceRepoWithEnv(initialDeviceRepoCtx, wg, conf.KafkaUrl, mongoUrl, conf.PermissionsUrl, map[string]string{
		"DISABLE_RIGHTS_HANDLING": "true",
		"SECURITY_IMPL":           "permissions-search",
	})
	if err != nil {
		t.Error(err)
		return
	}
	initialDeviceRepoUrl := "http://" + repoIp + ":8080"
	conf.DeviceRepoUrl = initialDeviceRepoUrl

	deviceOwners := []string{Userid, Userid, Userid, Userid, SecendOwnerTokenUser, SecendOwnerTokenUser, SecendOwnerTokenUser, SecendOwnerTokenUser}

	deviceIds := []string{} //will be filled in "init devices"

	//ints refer to deviceIds index
	expectedInitialAccess := map[string][]int{
		AdminTokenUser:       {0, 1, 2, 3, 4, 5, 6, 7},
		Userid:               {0, 1, 2, 3},
		SecendOwnerTokenUser: {4, 5, 6, 7},
	}
	//first update:
	//	remove admin role from 3 and 4
	//	allow UserId access to 6 and 7
	// 	allow SecendOwnerTokenUser access to 2 and 3
	//	remove UserId access to 2
	//	remove SecendOwnerTokenUser access to 6

	//key refers to deviceIds index
	firstUpdates := map[int]model.ResourceRightsBase{
		2: {
			UserRights: map[string]model.Right{
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{
				"admin": {Read: true, Write: true, Execute: true, Administrate: true},
			},
		},
		3: {
			UserRights: map[string]model.Right{
				Userid:               {Read: true, Write: true, Execute: true, Administrate: true},
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{},
		},
		4: {
			UserRights: map[string]model.Right{
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{},
		},
		6: {
			UserRights: map[string]model.Right{
				Userid: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{
				"admin": {Read: true, Write: true, Execute: true, Administrate: true},
			},
		},
		7: {
			UserRights: map[string]model.Right{
				Userid:               {Read: true, Write: true, Execute: true, Administrate: true},
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{
				"admin": {Read: true, Write: true, Execute: true, Administrate: true},
			},
		},
	}

	expectedAccessAfterFirstUpdate := map[string][]int{
		AdminTokenUser:       {0, 1, 2, 5, 6, 7},
		Userid:               {0, 1, 3, 6, 7},
		SecendOwnerTokenUser: {2, 3, 4, 5, 7},
	}

	//second update:
	//	remove admin role from 2 and 5
	//	allow UserId access to 4 and 5
	// 	allow SecendOwnerTokenUser access to 0 and 1
	//	remove UserId access to 0
	//	remove SecendOwnerTokenUser access to 4

	secondUpdates := map[int]model.ResourceRightsBase{
		0: {
			UserRights: map[string]model.Right{
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{
				"admin": {Read: true, Write: true, Execute: true, Administrate: true},
			},
		},
		1: {
			UserRights: map[string]model.Right{
				Userid:               {Read: true, Write: true, Execute: true, Administrate: true},
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{
				"admin": {Read: true, Write: true, Execute: true, Administrate: true},
			},
		},
		2: {
			UserRights: map[string]model.Right{
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{},
		},
		4: {
			UserRights: map[string]model.Right{
				Userid:               {Read: true, Write: true, Execute: true, Administrate: true},
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{},
		},
		5: {
			UserRights: map[string]model.Right{
				Userid:               {Read: true, Write: true, Execute: true, Administrate: true},
				SecendOwnerTokenUser: {Read: true, Write: true, Execute: true, Administrate: true},
			},
			GroupRights: map[string]model.Right{},
		},
	}

	expectedAccessAfterSecondUpdate := map[string][]int{
		AdminTokenUser:       {0, 1, 6, 7},
		Userid:               {1, 3, 4, 5, 6, 7},
		SecendOwnerTokenUser: {0, 1, 2, 3, 4, 5, 7},
	}

	getAdminTokenFromAccessMap := func(accessMap map[string][]int, deviceIndex int, intendedUpdate model.ResourceRightsBase) string {
		for userId, access := range accessMap {
			if slices.Contains(access, deviceIndex) && (intendedUpdate.UserRights[userId].Administrate || (userId == AdminTokenUser && intendedUpdate.GroupRights["admin"].Administrate)) {
				return userIdToUserToken[userId]
			}
		}
		return ""
	}

	t.Run("init devices", func(t *testing.T) {
		ctrl, err := controller.New(ctx, conf)
		if err != nil {
			t.Error(err)
			return
		}

		srv, err := api.Start(conf, ctrl)
		if err != nil {
			t.Fatal("ERROR: unable to start api", err)
		}
		defer srv.Shutdown(context.Background())

		time.Sleep(time.Second)

		protocol := models.Protocol{}
		t.Run("init protocols", func(t *testing.T) {
			resp, err := helper.Jwtpost(AdminToken, "http://localhost:"+conf.ServerPort+"/protocols?wait=true", models.Protocol{
				Name:             "p2",
				Handler:          "ph1",
				ProtocolSegments: []models.ProtocolSegment{{Name: "ps2"}},
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
				t.Error(err)
				return
			}
		})

		t.Run("init functions", func(t *testing.T) {
			resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/functions?wait=true", models.Function{Id: f1Id, Name: f1Id})
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				t.Fatal(resp.Status, resp.StatusCode, string(b))
			}
			resp.Body.Close()
		})

		t.Run("init aspects", func(t *testing.T) {
			resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/aspects?wait=true", models.Aspect{Id: a1Id, Name: a1Id})
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				t.Fatal(resp.Status, resp.StatusCode, string(b))
			}
			resp.Body.Close()
		})

		dt := models.DeviceType{}
		t.Run("init device-type", func(t *testing.T) {
			resp, err := helper.Jwtpost(AdminToken, "http://localhost:"+conf.ServerPort+"/device-types?wait=true", models.DeviceType{
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
				t.Error(err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				t.Fatal(resp.Status, resp.StatusCode, string(b))
			}

			err = json.NewDecoder(resp.Body).Decode(&dt)
			if err != nil {
				t.Error(err)
				return
			}

			if dt.Id == "" {
				t.Error(dt)
				return
			}
		})

		t.Run("create devices", func(t *testing.T) {
			for i, owner := range deviceOwners {
				token := userIdToUserToken[owner]
				resp, err := helper.Jwtpost(token, "http://localhost:"+conf.ServerPort+"/devices?wait=true", models.Device{
					LocalId:      fmt.Sprintf("foo-%v", i),
					Name:         fmt.Sprintf("foo-%v", i),
					DeviceTypeId: dt.Id,
				})
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					b, _ := io.ReadAll(resp.Body)
					t.Fatal(resp.Status, resp.StatusCode, string(b))
				}

				device := models.Device{}
				err = json.NewDecoder(resp.Body).Decode(&device)
				if err != nil {
					t.Error(err)
					return
				}
				deviceIds = append(deviceIds, device.Id)
			}
		})

	})

	t.Run("check initial rights", func(t *testing.T) {
		for user, allowedIds := range expectedInitialAccess {
			for idIndex, id := range deviceIds {
				resp, err := helper.Jwtget(userIdToUserToken[user], initialDeviceRepoUrl+"/devices/"+url.PathEscape(id))
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				if slices.Contains(allowedIds, idIndex) {
					if resp.StatusCode != http.StatusOK {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(resp.StatusCode, string(temp))
					}
				} else {
					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(resp.StatusCode, string(temp))
					}
				}
			}
		}
	})

	t.Run("update rights", func(t *testing.T) {
		for deviceIndex, update := range firstUpdates {
			token := getAdminTokenFromAccessMap(expectedInitialAccess, deviceIndex, update)
			deviceId := deviceIds[deviceIndex]

			timeout, _ := context.WithTimeout(ctx, time.Second*5)
			wait := donewait.AsyncWaitMultiple(timeout, []donewait.DoneMsg{
				{
					ResourceKind: "devices",
					ResourceId:   deviceId,
					Command:      "RIGHTS",
					Handler:      "github.com/SENERGY-Platform/permission-search",
				},
				{
					ResourceKind: "devices",
					ResourceId:   deviceId,
					Command:      "RIGHTS",
					Handler:      "github.com/SENERGY-Platform/device-repository",
				},
			}, nil)

			resp, err := helper.Jwtput(token, conf.PermissionsUrl+"/v3/administrate/rights/devices/"+url.PathEscape(deviceId), update)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				t.Fatal(deviceIndex, resp.Status, resp.StatusCode, string(b))
			}

			err = wait()
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	t.Run("check updated rights", func(t *testing.T) {
		for user, allowedIds := range expectedAccessAfterFirstUpdate {
			for idIndex, id := range deviceIds {
				resp, err := helper.Jwtget(userIdToUserToken[user], initialDeviceRepoUrl+"/devices/"+url.PathEscape(id))
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				if slices.Contains(allowedIds, idIndex) {
					if resp.StatusCode != http.StatusOK {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(resp.StatusCode, string(temp))
					}
				} else {
					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(resp.StatusCode, string(temp))
					}
				}
			}
		}
	})

	initialDeviceRepoCancel()

	time.Sleep(time.Minute)

	_, repoIp, err = docker.DeviceRepoWithEnv(ctx, wg, conf.KafkaUrl, mongoUrl, "-", map[string]string{
		"DISABLE_RIGHTS_HANDLING": "false",
		"SECURITY_IMPL":           "db",
	})
	if err != nil {
		t.Error(err)
		return
	}
	deviceRepoUrlWithInternalPermissions := "http://" + repoIp + ":8080"

	time.Sleep(10 * time.Second)

	t.Run("replay rights", func(t *testing.T) {
		err = permCmd(ctx, []string{"./app", "replay-permissions", "do", "devices"})
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(10 * time.Second)
	})

	t.Run("check rights after replay", func(t *testing.T) {
		for user, allowedIds := range expectedAccessAfterFirstUpdate {
			for idIndex, id := range deviceIds {
				resp, err := helper.Jwtget(userIdToUserToken[user], deviceRepoUrlWithInternalPermissions+"/devices/"+url.PathEscape(id))
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				if slices.Contains(allowedIds, idIndex) {
					if resp.StatusCode != http.StatusOK {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(user, idIndex, resp.StatusCode, string(temp))
					}
				} else {
					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(user, idIndex, resp.StatusCode, string(temp))
					}
				}
			}
		}
	})

	signal.Sub("", signal.Known.UpdateDone, func(value string, wg *sync.WaitGroup) {
		log.Printf("TEST-DEBUG: received update done %#v\n", value)
	})

	t.Run("update rights a second time", func(t *testing.T) {
		for deviceIndex, update := range secondUpdates {
			token := getAdminTokenFromAccessMap(expectedAccessAfterFirstUpdate, deviceIndex, update)
			if token == "" {
				t.Errorf("no valid token found for update %v %#v", deviceIndex, update)
				return
			}

			deviceId := deviceIds[deviceIndex]

			fmt.Println("update rights for", deviceIndex, deviceId)

			timeout, _ := context.WithTimeout(ctx, 5*time.Second)
			wait := donewait.AsyncWaitMultiple(timeout, []donewait.DoneMsg{
				{
					ResourceKind: "devices",
					ResourceId:   deviceId,
					Command:      "RIGHTS",
					Handler:      "github.com/SENERGY-Platform/permission-search",
				},
				{
					ResourceKind: "devices",
					ResourceId:   deviceId,
					Command:      "RIGHTS",
					Handler:      "github.com/SENERGY-Platform/device-repository",
				},
			}, nil)

			resp, err := helper.Jwtput(token, conf.PermissionsUrl+"/v3/administrate/rights/devices/"+url.PathEscape(deviceId), update)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				t.Error(deviceIndex, resp.Status, resp.StatusCode, string(b))
				return
			}

			err = wait()
			if err != nil {
				t.Error(deviceIndex, err)
				return
			}
		}
	})

	t.Run("check second updated rights", func(t *testing.T) {
		for user, allowedIds := range expectedAccessAfterSecondUpdate {
			for idIndex, id := range deviceIds {
				resp, err := helper.Jwtget(userIdToUserToken[user], deviceRepoUrlWithInternalPermissions+"/devices/"+url.PathEscape(id))
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()
				if slices.Contains(allowedIds, idIndex) {
					if resp.StatusCode != http.StatusOK {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(user, idIndex, resp.StatusCode, string(temp))
					}
				} else {
					if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusForbidden {
						temp, _ := io.ReadAll(resp.Body)
						t.Error(user, idIndex, resp.StatusCode, string(temp))
					}
				}
			}
		}
	})

}

const Userid = "testOwner"
const Userjwt = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOGM0N2E4OC0yYzc5LTQyMGYtODEwNC02NWJkOWViYmU0MWUiLCJleHAiOjE1NDY1MDcyMzMsIm5iZiI6MCwiaWF0IjoxNTQ2NTA3MTczLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwMDEvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJ0ZXN0T3duZXIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmcm9udGVuZCIsIm5vbmNlIjoiOTJjNDNjOTUtNzViMC00NmNmLTgwYWUtNDVkZDk3M2I0YjdmIiwiYXV0aF90aW1lIjoxNTQ2NTA3MDA5LCJzZXNzaW9uX3N0YXRlIjoiNWRmOTI4ZjQtMDhmMC00ZWI5LTliNjAtM2EwYWUyMmVmYzczIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ1c2VyIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibWFzdGVyLXJlYWxtIjp7InJvbGVzIjpbInZpZXctcmVhbG0iLCJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsIm1hbmFnZS1pZGVudGl0eS1wcm92aWRlcnMiLCJpbXBlcnNvbmF0aW9uIiwiY3JlYXRlLWNsaWVudCIsIm1hbmFnZS11c2VycyIsInF1ZXJ5LXJlYWxtcyIsInZpZXctYXV0aG9yaXphdGlvbiIsInF1ZXJ5LWNsaWVudHMiLCJxdWVyeS11c2VycyIsIm1hbmFnZS1ldmVudHMiLCJtYW5hZ2UtcmVhbG0iLCJ2aWV3LWV2ZW50cyIsInZpZXctdXNlcnMiLCJ2aWV3LWNsaWVudHMiLCJtYW5hZ2UtYXV0aG9yaXphdGlvbiIsIm1hbmFnZS1jbGllbnRzIiwicXVlcnktZ3JvdXBzIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJyb2xlcyI6WyJ1c2VyIl19.ykpuOmlpzj75ecSI6cHbCATIeY4qpyut2hMc1a67Ycg`

const AdminTokenUser = "admin"
const AdminToken = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOGM0N2E4OC0yYzc5LTQyMGYtODEwNC02NWJkOWViYmU0MWUiLCJleHAiOjE1NDY1MDcyMzMsIm5iZiI6MCwiaWF0IjoxNTQ2NTA3MTczLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwMDEvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJhZG1pbiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImZyb250ZW5kIiwibm9uY2UiOiI5MmM0M2M5NS03NWIwLTQ2Y2YtODBhZS00NWRkOTczYjRiN2YiLCJhdXRoX3RpbWUiOjE1NDY1MDcwMDksInNlc3Npb25fc3RhdGUiOiI1ZGY5MjhmNC0wOGYwLTRlYjktOWI2MC0zYTBhZTIyZWZjNzMiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbInVzZXIiLCJhZG1pbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7Im1hc3Rlci1yZWFsbSI6eyJyb2xlcyI6WyJ2aWV3LXJlYWxtIiwidmlldy1pZGVudGl0eS1wcm92aWRlcnMiLCJtYW5hZ2UtaWRlbnRpdHktcHJvdmlkZXJzIiwiaW1wZXJzb25hdGlvbiIsImNyZWF0ZS1jbGllbnQiLCJtYW5hZ2UtdXNlcnMiLCJxdWVyeS1yZWFsbXMiLCJ2aWV3LWF1dGhvcml6YXRpb24iLCJxdWVyeS1jbGllbnRzIiwicXVlcnktdXNlcnMiLCJtYW5hZ2UtZXZlbnRzIiwibWFuYWdlLXJlYWxtIiwidmlldy1ldmVudHMiLCJ2aWV3LXVzZXJzIiwidmlldy1jbGllbnRzIiwibWFuYWdlLWF1dGhvcml6YXRpb24iLCJtYW5hZ2UtY2xpZW50cyIsInF1ZXJ5LWdyb3VwcyJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwicm9sZXMiOlsidXNlciIsImFkbWluIl19.ggcFFFEsjwdfSzEFzmZt_m6W4IiSQub2FRhZVfWttDI`

const SecendOwnerTokenUser = "secondOwner"
const SecondOwnerToken = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOGM0N2E4OC0yYzc5LTQyMGYtODEwNC02NWJkOWViYmU0MWUiLCJleHAiOjE1NDY1MDcyMzMsIm5iZiI6MCwiaWF0IjoxNTQ2NTA3MTczLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwMDEvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJzZWNvbmRPd25lciIsInR5cCI6IkJlYXJlciIsImF6cCI6ImZyb250ZW5kIiwibm9uY2UiOiI5MmM0M2M5NS03NWIwLTQ2Y2YtODBhZS00NWRkOTczYjRiN2YiLCJhdXRoX3RpbWUiOjE1NDY1MDcwMDksInNlc3Npb25fc3RhdGUiOiI1ZGY5MjhmNC0wOGYwLTRlYjktOWI2MC0zYTBhZTIyZWZjNzMiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbInVzZXIiXX0sInJlc291cmNlX2FjY2VzcyI6eyJtYXN0ZXItcmVhbG0iOnsicm9sZXMiOlsidmlldy1yZWFsbSIsInZpZXctaWRlbnRpdHktcHJvdmlkZXJzIiwibWFuYWdlLWlkZW50aXR5LXByb3ZpZGVycyIsImltcGVyc29uYXRpb24iLCJjcmVhdGUtY2xpZW50IiwibWFuYWdlLXVzZXJzIiwicXVlcnktcmVhbG1zIiwidmlldy1hdXRob3JpemF0aW9uIiwicXVlcnktY2xpZW50cyIsInF1ZXJ5LXVzZXJzIiwibWFuYWdlLWV2ZW50cyIsIm1hbmFnZS1yZWFsbSIsInZpZXctZXZlbnRzIiwidmlldy11c2VycyIsInZpZXctY2xpZW50cyIsIm1hbmFnZS1hdXRob3JpemF0aW9uIiwibWFuYWdlLWNsaWVudHMiLCJxdWVyeS1ncm91cHMiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInJvbGVzIjpbInVzZXIiXX0.cq8YeUuR0jSsXCEzp634fTzNbGkq_B8KbVrwBPgceJ4`

var userIdToUserToken = map[string]string{
	Userid:               Userjwt,
	SecendOwnerTokenUser: SecondOwnerToken,
	AdminTokenUser:       AdminToken,
}
