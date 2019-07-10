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

package services

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
	"github.com/SENERGY-Platform/device-manager/lib/tests/servicemocks"
	"github.com/ory/dockertest"
	"log"
	"sync"
)

func New(cin config.Config) (publ *servicemocks.Publisher, cout config.Config, shutdown func(), err error) {
	cout = cin

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Println("Could not connect to docker: %s", err)
		return publ, cout, func() {}, err
	}

	listMux := sync.Mutex{}
	var globalError error
	closerList := []func(){}
	close := func(list []func()) {
		for i := len(list)/2 - 1; i >= 0; i-- {
			opp := len(list) - 1 - i
			list[i], list[opp] = list[opp], list[i]
		}
		for _, c := range list {
			if c != nil {
				c()
			}
		}
	}

	defer func() {
		if globalError != nil {
			close(closerList)
		}
	}()

	var elasticIp string

	closeZk, _, zkIp, err := Zookeeper(pool)
	listMux.Lock()
	closerList = append(closerList, closeZk)
	listMux.Unlock()
	if err != nil {
		globalError = err
		return
	}
	cout.ZookeeperUrl = zkIp + ":2181"

	//kafka
	closeKafka, err := Kafka(pool, cout.ZookeeperUrl)
	listMux.Lock()
	closerList = append(closerList, closeKafka)
	listMux.Unlock()
	if err != nil {
		globalError = err
		return
	}

	closeElastic, _, ip, err := Elasticsearch(pool)
	elasticIp = ip
	listMux.Lock()
	closerList = append(closerList, closeElastic)
	listMux.Unlock()
	if err != nil {
		globalError = err
		return
	}

	//permsearch
	closePerm, _, permIp, err := PermSearch(pool, cout.ZookeeperUrl, elasticIp)
	listMux.Lock()
	closerList = append(closerList, closePerm)
	listMux.Unlock()
	if err != nil {
		globalError = err
		return
	}
	cout.PermissionsUrl = "http://" + permIp + ":8080"

	closerMongo, _, mongoIp, err := MongoTestServer(pool)
	listMux.Lock()
	closerList = append(closerList, closerMongo)
	listMux.Unlock()
	if err != nil {
		globalError = err
		return
	}

	closerDeviceRepo, _, repoIp, err := DeviceRepo(pool, mongoIp, cout.ZookeeperUrl, permIp)
	listMux.Lock()
	closerList = append(closerList, closerDeviceRepo)
	listMux.Unlock()
	if err != nil {
		globalError = err
		return
	}
	cout.DeviceRepoUrl = "http://" + repoIp + ":8080"

	kafkapubl, err := publisher.New(cout)
	if err != nil {
		close(closerList)
		return publ, cout, shutdown, err
	}

	publ = servicemocks.NewPublisher()
	publ.Subscribe(servicemocks.DtTopic, func(msg []byte) {
		cmd := publisher.DeviceTypeCommand{}
		json.Unmarshal(msg, &cmd)
		kafkapubl.PublishDeviceTypeCommand(cmd)
	})
	publ.Subscribe(servicemocks.ProtocolTopic, func(msg []byte) {
		cmd := publisher.ProtocolCommand{}
		json.Unmarshal(msg, &cmd)
		kafkapubl.PublishProtocolCommand(cmd)
	})

	semantic := servicemocks.NewSemanticRepo(publ)
	cout.SemanticRepoUrl = semantic.Url()

	return publ, cout, func() {
		semantic.Stop()
		close(closerList)
	}, nil
}
