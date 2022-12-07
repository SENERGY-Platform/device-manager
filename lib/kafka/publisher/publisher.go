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

package publisher

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/util"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"os"
)

type Publisher struct {
	config          config.Config
	devicetypes     *kafka.Writer
	devicegroups    *kafka.Writer
	protocols       *kafka.Writer
	devices         *kafka.Writer
	hubs            *kafka.Writer
	concepts        *kafka.Writer
	characteristics *kafka.Writer
	aspects         *kafka.Writer
	functions       *kafka.Writer
	deviceclasses   *kafka.Writer
	locations       *kafka.Writer
	permissions     *kafka.Writer
}

func New(conf config.Config, ctx context.Context) (*Publisher, error) {
	log.Println("ensure kafka topics")
	err := util.InitTopic(
		conf.KafkaUrl,
		conf.DeviceTypeTopic,
		conf.DeviceGroupTopic,
		conf.ProtocolTopic,
		conf.DeviceTopic,
		conf.HubTopic,
		conf.ConceptTopic,
		conf.CharacteristicTopic,
		conf.AspectTopic,
		conf.FunctionTopic,
		conf.DeviceClassTopic,
		conf.LocationTopic,
		conf.PermissionsTopic)
	if err != nil {
		return nil, err
	}
	broker, err := util.GetBroker(conf.KafkaUrl)
	if err != nil {
		return nil, err
	}
	if len(broker) == 0 {
		return nil, errors.New("missing kafka broker")
	}
	log.Println("Produce to ", conf.DeviceTypeTopic, conf.ProtocolTopic, conf.DeviceTopic, conf.HubTopic, conf.ConceptTopic, conf.CharacteristicTopic, conf.LocationTopic)
	devicetypes := getProducer(ctx, broker, conf.DeviceTypeTopic, conf.LogLevel == "DEBUG")
	devicegroups := getProducer(ctx, broker, conf.DeviceGroupTopic, conf.LogLevel == "DEBUG")
	devices := getProducer(ctx, broker, conf.DeviceTopic, conf.LogLevel == "DEBUG")
	hubs := getProducer(ctx, broker, conf.HubTopic, conf.LogLevel == "DEBUG")
	protocol := getProducer(ctx, broker, conf.ProtocolTopic, conf.LogLevel == "DEBUG")
	concepts := getProducer(ctx, broker, conf.ConceptTopic, conf.LogLevel == "DEBUG")
	characteristics := getProducer(ctx, broker, conf.CharacteristicTopic, conf.LogLevel == "DEBUG")
	aspect := getProducer(ctx, broker, conf.AspectTopic, conf.LogLevel == "DEBUG")
	function := getProducer(ctx, broker, conf.FunctionTopic, conf.LogLevel == "DEBUG")
	deviceclass := getProducer(ctx, broker, conf.DeviceClassTopic, conf.LogLevel == "DEBUG")
	location := getProducer(ctx, broker, conf.LocationTopic, conf.LogLevel == "DEBUG")
	permissions := getProducer(ctx, broker, conf.PermissionsTopic, conf.LogLevel == "DEBUG")
	return &Publisher{
		config:          conf,
		devicetypes:     devicetypes,
		devicegroups:    devicegroups,
		protocols:       protocol,
		devices:         devices,
		hubs:            hubs,
		concepts:        concepts,
		characteristics: characteristics,
		aspects:         aspect,
		functions:       function,
		deviceclasses:   deviceclass,
		locations:       location,
		permissions:     permissions,
	}, nil
}

func getProducer(ctx context.Context, broker []string, topic string, debug bool) (writer *kafka.Writer) {
	var logger *log.Logger
	if debug {
		logger = log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0)
	} else {
		logger = log.New(io.Discard, "", 0)
	}
	writer = &kafka.Writer{
		Addr:        kafka.TCP(broker...),
		Topic:       topic,
		MaxAttempts: 10,
		Logger:      logger,
		BatchSize:   1,
		Balancer:    &kafka.Hash{},
	}
	go func() {
		<-ctx.Done()
		err := writer.Close()
		if err != nil {
			log.Println("ERROR: unable to close producer for", topic, err)
		}
	}()
	return writer
}
