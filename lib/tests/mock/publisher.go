package mock

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
)

const DtTopic = "device-type"

type Publisher struct {
	listener map[string][]func(msg []byte)
}

func NewPublisher() *Publisher {
	return &Publisher{listener: map[string][]func(msg []byte){}}
}

func (this *Publisher) PublishDeviceType(device model.DeviceType, userId string) (err error) {
	cmd := publisher.DeviceTypeCommand{Command: "PUT", Id: device.Id, DeviceType: device, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DtTopic, message)
}

func (this *Publisher) PublishDeviceDelete(id string, userId string) error {
	cmd := publisher.DeviceTypeCommand{Command: "DELETE", Id: id, Owner: userId}
	message, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.send(DtTopic, message)
}

func (this *Publisher) send(topic string, msg []byte) error {
	for _, listener := range this.listener[topic] {
		go listener(msg)
	}
	return nil
}

func (this *Publisher) Subscribe(topic string, f func(msg []byte)) {
	this.listener[topic] = append(this.listener[topic], f)
}
