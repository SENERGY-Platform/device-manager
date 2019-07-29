package model

import "github.com/google/uuid"

func (variable *Characteristic) GenerateId() {
	variable.Id = "urn:infai:ses:categoryvariable:" + uuid.New().String()
	for i, v := range variable.SubCharacteristics {
		v.GenerateId()
		variable.SubCharacteristics[i] = v
	}
}

func (class *DeviceClass) GenerateId() {
	class.Id = "urn:infai:ses:device-class:" + uuid.New().String()
}

func (function *Function) GenerateId() {
	function.Id = "urn:infai:ses:function:" + uuid.New().String()
}

func (aspect *Aspect) GenerateId() {
	aspect.Id = "urn:infai:ses:aspect:" + uuid.New().String()
}

func (category *Concept) GenerateId() {
	category.Id = "urn:infai:ses:category:" + uuid.New().String()
	for i, v := range category.Characteristics {
		v.GenerateId()
		category.Characteristics[i] = v
	}
}

func (device *Device) GenerateId() {
	device.Id = "urn:infai:ses:device:" + uuid.New().String()
}

func (deviceType *DeviceType) GenerateId() {
	deviceType.Id = "urn:infai:ses:device-type:" + uuid.New().String()
	for i, service := range deviceType.Services {
		service.GenerateId()
		deviceType.Services[i] = service
	}
	if deviceType.DeviceClass.Id == "" {
		deviceType.DeviceClass.GenerateId()
	}
}

func (service *Service) GenerateId() {
	service.Id = "urn:infai:ses:service:" + uuid.New().String()
	for i, function := range service.Functions {
		if function.Id == "" {
			function.GenerateId()
			service.Functions[i] = function
		}
	}
	for i, aspect := range service.Aspects {
		if aspect.Id == "" {
			aspect.GenerateId()
			service.Aspects[i] = aspect
		}
	}
	for i, content := range service.Inputs {
		content.GenerateId()
		service.Inputs[i] = content
	}
	for i, content := range service.Outputs {
		content.GenerateId()
		service.Outputs[i] = content
	}
}

func (hub *Hub) GenerateId() {
	hub.Id = "urn:infai:ses:hub:" + uuid.New().String()
}

func (protocol *Protocol) GenerateId() {
	protocol.Id = "urn:infai:ses:protocol:" + uuid.New().String()
	for i, segment := range protocol.ProtocolSegments {
		segment.GenerateId()
		protocol.ProtocolSegments[i] = segment
	}
}

func (segment *ProtocolSegment) GenerateId() {
	segment.Id = "urn:infai:ses:segment:" + uuid.New().String()
}

func (content *Content) GenerateId() {
	content.Id = "urn:infai:ses:content:" + uuid.New().String()
	content.ContentVariable.GenerateId()
}

func (variable *ContentVariable) GenerateId() {
	variable.Id = "urn:infai:ses:contentvariable:" + uuid.New().String()
	for i, v := range variable.SubContentVariables {
		v.GenerateId()
		variable.SubContentVariables[i] = v
	}
}
