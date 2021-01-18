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
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
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
}

func New(conf config.Config) (*Publisher, error) {
	log.Println("ensure kafka topics")
	err := InitTopicWithConfig(
		conf.ZookeeperUrl,
		1,
		1,
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
		conf.LocationTopic)
	if err != nil {
		return nil, err
	}
	broker, err := GetBroker(conf.ZookeeperUrl)
	if err != nil {
		return nil, err
	}
	if len(broker) == 0 {
		return nil, errors.New("missing kafka broker")
	}
	log.Println("Produce to ", conf.DeviceTypeTopic, conf.ProtocolTopic, conf.DeviceTopic, conf.HubTopic, conf.ConceptTopic, conf.CharacteristicTopic, conf.LocationTopic)
	devicetypes, err := getProducer(broker, conf.DeviceTypeTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	devicegroups, err := getProducer(broker, conf.DeviceGroupTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	devices, err := getProducer(broker, conf.DeviceTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	hubs, err := getProducer(broker, conf.HubTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	protocol, err := getProducer(broker, conf.ProtocolTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	concepts, err := getProducer(broker, conf.ConceptTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	characteristics, err := getProducer(broker, conf.CharacteristicTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	aspect, err := getProducer(broker, conf.AspectTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	function, err := getProducer(broker, conf.FunctionTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	deviceclass, err := getProducer(broker, conf.DeviceClassTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	location, err := getProducer(broker, conf.LocationTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func getProducer(broker []string, topic string, debug bool) (writer *kafka.Writer, err error) {
	var logger *log.Logger
	if debug {
		logger = log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0)
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:     broker,
		Topic:       topic,
		MaxAttempts: 10,
		Logger:      logger,
	})
	return writer, err
}
