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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/listener"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/publisher"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/util"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/tests/docker"
	"github.com/SENERGY-Platform/device-manager/lib/tests/servicemocks"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
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

	mockPublisher := servicemocks.NewPublisher()
	repo := servicemocks.NewDeviceRepo(mockPublisher)
	conf.DeviceRepoUrl = repo.Url()
	semantic := servicemocks.NewSemanticRepo(mockPublisher)
	conf.SemanticRepoUrl = semantic.Url()

	ctrl, err := controller.New(ctx, conf)
	if err != nil {
		t.Error(err)
		return
	}

	broker, err := util.GetBroker(conf.KafkaUrl)
	if err != nil {
		t.Error(err)
		return
	}
	if len(broker) == 0 {
		t.Error("missing kafka broker")
		return
	}

	t.Run("create devices", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			id := strconv.Itoa(i)
			device := model.Device{
				Id:           id,
				LocalId:      id,
				Name:         id + "_name",
				Attributes:   nil,
				DeviceTypeId: "test_dt",
			}
			_, err, _ = ctrl.PublishDeviceUpdate(user1a, id, device)
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 20; i < 40; i++ {
			id := strconv.Itoa(i)
			device := model.Device{
				Id:           id,
				LocalId:      id,
				Name:         id + "_name",
				Attributes:   nil,
				DeviceTypeId: "test_dt",
			}
			log.Println("test create device", id)
			_, err, _ = ctrl.PublishDeviceUpdate(user2a, id, device)
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	time.Sleep(10 * time.Second)

	t.Run("change permissions", func(t *testing.T) {
		permissions := &kafka.Writer{
			Addr:        kafka.TCP(broker...),
			Topic:       conf.PermissionsTopic,
			MaxAttempts: 10,
			Logger:      log.New(os.Stdout, "[TEST-KAFKA-PRODUCER] ", 0),
		}
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i)
			err = setPermission(permissions, conf, user2.GetUserId(), id, "rwxa")
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 20; i < 30; i++ {
			id := strconv.Itoa(i)
			err = setPermission(permissions, conf, user1.GetUserId(), id, "rwxa")
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 5; i < 10; i++ {
			id := strconv.Itoa(i)
			err = setPermission(permissions, conf, user1.GetUserId(), id, "rx")
			if err != nil {
				t.Error(err)
				return
			}
		}
		for i := 25; i < 30; i++ {
			id := strconv.Itoa(i)
			err = setPermission(permissions, conf, user2.GetUserId(), id, "rx")
			if err != nil {
				t.Error(err)
				return
			}
		}
	})

	time.Sleep(60 * time.Second)

	t.Run("check user1 before delete", checkUserDevices(conf, user1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}))
	t.Run("check user2 before delete", checkUserDevices(conf, user2, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39}))

	t.Run("delete user1", func(t *testing.T) {
		users := &kafka.Writer{
			Addr:        kafka.TCP(broker...),
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

	time.Sleep(60 * time.Second)

	t.Run("check user1 after delete", checkUserDevices(conf, user1, []int{}))
	t.Run("check user2 after delete", checkUserDevices(conf, user2, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 20, 21, 22, 23, 24, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39}))

}

func checkUserDevices(conf config.Config, token auth.Token, expectedDeviceIdsAsInt []int) func(t *testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest("GET", conf.PermissionsUrl+"/v3/resources/devices?rights=r&limit=100", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token.Token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			resp.Body.Close()
			log.Println("DEBUG: PermissionCheck()", buf.String())
			err = errors.New("access denied")
			t.Error(err)
			return
		}

		devices := []map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&devices)
		if err != nil {
			t.Error(err)
			return
		}
		actualIds := []string{}
		for _, device := range devices {
			id, ok := device["id"].(string)
			if !ok {
				t.Error("expect device id to be string", device)
				return
			}
			actualIds = append(actualIds, id)
		}
		sort.Strings(actualIds)

		expectedIds := []string{}
		for _, intId := range expectedDeviceIdsAsInt {
			expectedIds = append(expectedIds, strconv.Itoa(intId))
		}
		sort.Strings(expectedIds)

		if !reflect.DeepEqual(actualIds, expectedIds) {
			t.Error(actualIds, expectedIds)
			return
		}
	}
}

func setPermission(permissions *kafka.Writer, conf config.Config, userId string, deviceId string, right string) error {
	cmd := publisher.PermCommandMsg{
		Command:  "PUT",
		Kind:     conf.DeviceTopic,
		Resource: deviceId,
		User:     userId,
		Right:    right,
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	err = permissions.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(userId + "_" + conf.DeviceTopic + "_" + deviceId),
			Value: message,
			Time:  time.Now(),
		},
	)
	return err
}
