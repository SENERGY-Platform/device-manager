/*
 * Copyright 2021 InfAI (CC SES)
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
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/listener"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/tests/docker"
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/permissions-v2/pkg/client"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestUserDelete(t *testing.T) {
	conf, err := config.Load("./../../config.json")
	if err != nil {
		t.Fatal("ERROR: unable to load config", err)
	}
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user1, err := auth.CreateToken("test", "user1")
	if err != nil {
		t.Error(err)
		return
	}
	user2, err := auth.CreateToken("test", "user2")
	if err != nil {
		t.Error(err)
		return
	}
	user1a, err := auth.CreateTokenWithRoles("test", "user1", []string{"admin"})
	if err != nil {
		t.Error(err)
		return
	}
	user2a, err := auth.CreateTokenWithRoles("test", "user2", []string{"admin"})
	if err != nil {
		t.Error(err)
		return
	}

	conf, err = docker.CreateTestEnv(ctx, wg, conf)
	if err != nil {
		t.Error(err)
		return
	}

	//to ensure that pagination is used
	oldBatchSize := com.ResourcesEffectedByUserDelete_BATCH_SIZE
	com.ResourcesEffectedByUserDelete_BATCH_SIZE = 5
	defer func() {
		com.ResourcesEffectedByUserDelete_BATCH_SIZE = oldBatchSize
	}()

	time.Sleep(10 * time.Second)

	conf.Debug = true

	ctrl, err := controller.New(ctx, conf)
	if err != nil {
		t.Error(err)
		return
	}

	cache := &map[string]client.ResourcePermissions{}

	dt := models.DeviceType{}
	t.Run("create device-type", func(t *testing.T) {
		protocol, err, _ := ctrl.PublishProtocolCreate(user1a, models.Protocol{
			Name:             "p2",
			Handler:          "ph1",
			ProtocolSegments: []models.ProtocolSegment{{Name: "ps2"}},
		}, model.ProtocolUpdateOptions{Wait: true})
		if err != nil {
			t.Error(err)
			return
		}

		dt, err, _ = ctrl.PublishDeviceTypeCreate(user1a, models.DeviceType{
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
								Name: "v1name",
								Type: models.String,
							},
						},
					},

					ProtocolId: protocol.Id,
				},
			},
		}, model.DeviceTypeUpdateOptions{Wait: true})
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(2 * time.Second)
	})

	t.Run("create devices", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			id := strconv.Itoa(i)
			device := models.Device{
				Id:           id,
				LocalId:      id,
				Name:         id + "_name",
				Attributes:   nil,
				DeviceTypeId: dt.Id,
			}
			_, err, _ = ctrl.PublishDeviceUpdate(user1a, id, device, model.DeviceUpdateOptions{})
			if err != nil {
				t.Error(err)
				return
			}
			(*cache)[conf.DeviceTopic+"."+device.Id] = client.ResourcePermissions{
				UserPermissions: map[string]client.PermissionsMap{user1a.GetUserId(): {
					Read:         true,
					Write:        true,
					Execute:      true,
					Administrate: true,
				}},
				RolePermissions: map[string]client.PermissionsMap{"admin": {Read: true, Write: true, Execute: true, Administrate: true}},
			}
		}
		for i := 20; i < 40; i++ {
			id := strconv.Itoa(i)
			device := models.Device{
				Id:           id,
				LocalId:      id,
				Name:         id + "_name",
				Attributes:   nil,
				DeviceTypeId: dt.Id,
			}
			log.Println("test create device", id)
			_, err, _ = ctrl.PublishDeviceUpdate(user2a, id, device, model.DeviceUpdateOptions{})
			if err != nil {
				t.Error(err)
				return
			}
			(*cache)[conf.DeviceTopic+"."+device.Id] = client.ResourcePermissions{
				UserPermissions: map[string]client.PermissionsMap{user2a.GetUserId(): {
					Read:         true,
					Write:        true,
					Execute:      true,
					Administrate: true,
				}},
				RolePermissions: map[string]client.PermissionsMap{"admin": {Read: true, Write: true, Execute: true, Administrate: true}},
			}
		}
	})

	t.Run("change permissions", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i)
			err = setPermission(ctrl.GetCom(), user2.GetUserId(), conf.DeviceTopic, id, "rwxa", cache)
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 20; i < 30; i++ {
			id := strconv.Itoa(i)
			err = setPermission(ctrl.GetCom(), user1.GetUserId(), conf.DeviceTopic, id, "rwxa", cache)
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 5; i < 10; i++ {
			id := strconv.Itoa(i)
			err = setPermission(ctrl.GetCom(), user1.GetUserId(), conf.DeviceTopic, id, "rx", cache)
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 25; i < 30; i++ {
			id := strconv.Itoa(i)
			err = setPermission(ctrl.GetCom(), user2.GetUserId(), conf.DeviceTopic, id, "rx", cache)
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	time.Sleep(10 * time.Second)

	t.Run("check user1 before delete", checkUserDevices(conf, user1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}))
	t.Run("check user2 before delete", checkUserDevices(conf, user2, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39}))

	t.Run("delete user1", func(t *testing.T) {
		users := &kafka.Writer{
			Addr:        kafka.TCP(conf.KafkaUrl),
			Topic:       conf.UserTopic,
			MaxAttempts: 10,
			Logger:      log.New(os.Stdout, "[TEST-KAFKA-PRODUCER] ", 0),
		}
		cmd := listener.UserCommandMsg{
			Command: "DELETE",
			Id:      user1.GetUserId(),
		}
		message, err := json.Marshal(cmd)
		if err != nil {
			t.Error(err)
			return
		}
		err = users.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte(user1.GetUserId()),
				Value: message,
				Time:  time.Now(),
			},
		)
	})

	time.Sleep(10 * time.Second)

	t.Run("check user1 after delete", checkUserDevices(conf, user1, []int{}))
	t.Run("check user2 after delete", checkUserDevices(conf, user2, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39}))

}

func checkUserDevices(conf config.Config, token auth.Token, expectedDeviceIdsAsInt []int) func(t *testing.T) {
	return func(t *testing.T) {
		devices, err, _ := devicerepo.NewClient(conf.DeviceRepoUrl, nil).ListDevices(token.Jwt(), devicerepo.DeviceListOptions{Limit: 100})
		if err != nil {
			t.Error(err)
			return
		}
		actualIds := []string{}
		for _, device := range devices {
			actualIds = append(actualIds, device.Id)
		}
		sort.Strings(actualIds)

		expectedIds := []string{}
		for _, intId := range expectedDeviceIdsAsInt {
			expectedIds = append(expectedIds, strconv.Itoa(intId))
		}
		sort.Strings(expectedIds)
		if !reflect.DeepEqual(actualIds, expectedIds) {
			t.Errorf("\na=%#v\ne=%#v\n", actualIds, expectedIds)
			return
		}
	}
}

func setPermission(com controller.Com, userId string, kind string, id string, right string, cache *map[string]client.ResourcePermissions) error {
	userRight := client.PermissionsMap{
		Read:         false,
		Write:        false,
		Execute:      false,
		Administrate: false,
	}
	for _, r := range right {
		switch r {
		case 'r':
			userRight.Read = true
		case 'w':
			userRight.Write = true
		case 'a':
			userRight.Administrate = true
		case 'x':
			userRight.Execute = true
		default:
			return errors.New("unknown right in " + right)
		}
	}
	cacheKey := kind + "." + id
	msg, ok := (*cache)[cacheKey]
	if !ok {
		msg = client.ResourcePermissions{
			UserPermissions: map[string]client.PermissionsMap{},
			RolePermissions: map[string]client.PermissionsMap{"admin": {Read: true, Write: true, Execute: true, Administrate: true}},
		}
	}
	msg.UserPermissions[userId] = userRight
	(*cache)[cacheKey] = msg

	_, err, _ := com.SetPermission(client.InternalAdminToken, kind, id, msg)
	return err
}
