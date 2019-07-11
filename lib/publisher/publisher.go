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
	config      config.Config
	devicetypes *kafka.Writer
	protocols   *kafka.Writer
}

func New(conf config.Config) (*Publisher, error) {
	broker, err := GetBroker(conf.ZookeeperUrl)
	if err != nil {
		return nil, err
	}
	if len(broker) == 0 {
		return nil, errors.New("missing kafka broker")
	}
	log.Println("Produce to ", conf.DeviceTypeTopic, conf.ProtocolTopic)
	devicetypes, err := getProducer(broker, conf.DeviceTypeTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	protocol, err := getProducer(broker, conf.ProtocolTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	return &Publisher{config: conf, devicetypes: devicetypes, protocols: protocol}, nil
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
