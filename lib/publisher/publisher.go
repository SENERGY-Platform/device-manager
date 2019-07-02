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
}

func New(conf config.Config) (*Publisher, error) {
	broker, err := GetBroker(conf.ZookeeperUrl)
	if err != nil {
		return nil, err
	}
	if len(broker) == 0 {
		return nil, errors.New("missing kafka broker")
	}
	devicetypes, err := getProducer(broker, conf.DeviceTypeTopic, conf.LogLevel == "DEBUG")
	if err != nil {
		return nil, err
	}
	return &Publisher{config: conf, devicetypes: devicetypes}, nil
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
