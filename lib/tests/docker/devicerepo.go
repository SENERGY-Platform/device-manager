/*
 * Copyright 2020 InfAI (CC SES)
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

package docker

import (
	"context"
	"github.com/ory/dockertest/v3"
	"log"
	"net/http"
	"sync"
	"time"
)

func DeviceRepoWithDependencies(basectx context.Context, wg *sync.WaitGroup) (repoUrl string, searchUrl string, kafkaUrl string, err error) {
	ctx, cancel := context.WithCancel(basectx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", "", "", err
	}

	_, zkIp, err := Zookeeper(pool, ctx, wg)
	if err != nil {
		return repoUrl, searchUrl, kafkaUrl, err
	}
	zookeeperUrl := zkIp + ":2181"

	kafkaUrl, err = Kafka(pool, ctx, wg, zookeeperUrl)
	if err != nil {
		return repoUrl, searchUrl, kafkaUrl, err
	}

	time.Sleep(1 * time.Second)

	_, elasticIp, err := Elasticsearch(pool, ctx, wg)
	if err != nil {
		return repoUrl, searchUrl, kafkaUrl, err
	}

	_, permIp, err := PermSearch(pool, ctx, wg, kafkaUrl, elasticIp)
	if err != nil {
		return repoUrl, searchUrl, kafkaUrl, err
	}
	searchUrl = "http://" + permIp + ":8080"

	_, mongoIp, err := MongoDB(pool, ctx, wg)
	if err != nil {
		return repoUrl, searchUrl, kafkaUrl, err
	}

	_, repoIp, err := DeviceRepo(pool, ctx, wg, kafkaUrl, "mongodb://"+mongoIp+":27017", searchUrl)
	if err != nil {
		return repoUrl, searchUrl, kafkaUrl, err
	}
	repoUrl = "http://" + repoIp + ":8080"

	return repoUrl, searchUrl, kafkaUrl, err
}

func DeviceRepo(pool *dockertest.Pool, ctx context.Context, wg *sync.WaitGroup, kafkaUrl string, mongoUrl string, permsearch string) (hostPort string, ipAddress string, err error) {
	log.Println("start device-repository")
	container, err := pool.Run("ghcr.io/senergy-platform/device-repository", "dev", []string{
		"KAFKA_URL=" + kafkaUrl,
		"PERMISSIONS_URL=" + permsearch,
		"MONGO_URL=" + mongoUrl,
	})
	if err != nil {
		return "", "", err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("DEBUG: remove container " + container.Container.Name)
		container.Close()
	}()
	go Dockerlog(pool, ctx, container, "DEVICE-REPOSITORY")
	hostPort = container.GetPort("8080/tcp")
	err = pool.Retry(func() error {
		log.Println("try device-repo connection...")
		_, err := http.Get("http://localhost:" + hostPort)
		if err != nil {
			log.Println(err)
		}
		return err
	})
	return hostPort, container.Container.NetworkSettings.IPAddress, err
}
