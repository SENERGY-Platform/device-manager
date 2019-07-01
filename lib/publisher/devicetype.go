package publisher

import (
	"context"
	"encoding/json"
	"github.com/SmartEnergyPlatform/device-manager/lib/model"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type DeviceTypeCommand struct {
	Command    string           `json:"command"`
	Id         string           `json:"id"`
	Owner      string           `json:"owner"`
	DeviceType model.DeviceType `json:"device"`
}

func (this *Publisher) PublishDeviceType(device model.DeviceType, userId string) (err error) {
	cmd := DeviceTypeCommand{Command: "PUT", Id: device.Id, DeviceType: device, Owner: userId}
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce", cmd)
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.devicetypes.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(device.Id),
			Value: message,
			Time:  time.Now(),
		},
	)
}

func (this *Publisher) PublishDeviceDelete(id string, userId string) error {
	cmd := DeviceTypeCommand{Command: "DELETE", Id: id, Owner: userId}
	if this.config.LogLevel == "DEBUG" {
		log.Println("DEBUG: produce", cmd)
	}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.devicetypes.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(id),
			Value: message,
			Time:  time.Now(),
		},
	)
}
