package controller

import (
	"github.com/SmartEnergyPlatform/device-manager/lib/config"
	"github.com/SmartEnergyPlatform/device-manager/lib/controller/com"
	"github.com/SmartEnergyPlatform/device-manager/lib/model"
	"github.com/SmartEnergyPlatform/device-manager/lib/publisher"
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
