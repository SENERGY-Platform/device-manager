package controller

import (
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/controller/com"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/publisher"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

type Controller struct {
	publisher Publisher
	com       Com
}

func New(conf config.Config) (*Controller, error) {
	publ, err := publisher.New(conf)
	if err != nil {
		return &Controller{}, err
	}
	return &Controller{com: com.New(conf), publisher: publ}, nil
}

func NewWithPublisher(conf config.Config, publisher Publisher) (*Controller, error) {
	return &Controller{com: com.New(conf), publisher: publisher}, nil
}

type Publisher interface {
	PublishDeviceType(device model.DeviceType, userID string) (err error)
	PublishDeviceDelete(id string, userID string) error
}

type Com interface {
	ValidateDeviceType(jwt jwt_http_router.Jwt, dt model.DeviceType) (err error, code int)
	PermissionCheckForDeviceType(jwt jwt_http_router.Jwt, id string, permission string) (err error, code int) //permission = "w" | "r" | "x" | "a"
	GetTechnicalDeviceType(jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int)
	GetSemanticDeviceType(jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int)
}
