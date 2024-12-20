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

package docker

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func CreateTestEnv(baseCtx context.Context, wg *sync.WaitGroup, startConfig config.Config) (config config.Config, err error) {
	config = startConfig
	ctx, cancel := context.WithCancel(baseCtx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	whPort, err := getFreePort()
	if err != nil {
		log.Println("unable to find free port", err)
		return config, err
	}
	config.ServerPort = strconv.Itoa(whPort)

	_, zkIp, err := Zookeeper(ctx, wg)
	if err != nil {
		return config, err
	}
	zookeeperUrl := zkIp + ":2181"

	config.KafkaUrl, err = Kafka(ctx, wg, zookeeperUrl)
	if err != nil {
		return config, err
	}

	_, mongoIp, err := MongoDB(ctx, wg)
	if err != nil {
		return config, err
	}
	mongoUrl := "mongodb://" + mongoIp + ":27017"

	_, permV2Ip, err := PermissionsV2(ctx, wg, mongoUrl, config.KafkaUrl)
	if err != nil {
		return config, err
	}
	config.PermissionsV2Url = "http://" + permV2Ip + ":8080"

	_, repoIp, err := DeviceRepo(ctx, wg, config.KafkaUrl, mongoUrl, config.PermissionsV2Url)
	if err != nil {
		return config, err
	}
	config.DeviceRepoUrl = "http://" + repoIp + ":8080"

	time.Sleep(10 * time.Second)
	return
}

func Dockerlog(container testcontainers.Container, name string) error {
	container.FollowOutput(&LogWriter{logger: log.New(os.Stdout, "["+name+"] ", log.LstdFlags)})
	err := container.StartLogProducer(context.Background())
	if err != nil {
		return err
	}
	return nil
}

type LogWriter struct {
	logger *log.Logger
}

func (this *LogWriter) Accept(l testcontainers.Log) {
	this.Write(l.Content)
}

func (this *LogWriter) Write(p []byte) (n int, err error) {
	this.logger.Print(string(p))
	return len(p), nil
}

func Forward(ctx context.Context, fromPort int, toAddr string) error {
	log.Println("forward", fromPort, "to", toAddr)
	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", fromPort))
	if err != nil {
		return err
	}
	go func() {
		defer log.Println("closed forward incoming")
		<-ctx.Done()
		incoming.Close()
	}()
	go func() {
		for {
			client, err := incoming.Accept()
			if err != nil {
				log.Println("FORWARD ERROR:", err)
				return
			}
			go handleForwardClient(client, toAddr)
		}
	}()
	return nil
}

func handleForwardClient(client net.Conn, addr string) {
	//log.Println("new forward client")
	target, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("FORWARD ERROR:", err)
		return
	}
	go func() {
		defer target.Close()
		defer client.Close()
		io.Copy(target, client)
	}()
	go func() {
		defer target.Close()
		defer client.Close()
		io.Copy(client, target)
	}()
}
