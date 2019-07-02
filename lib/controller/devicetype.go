package controller

import (
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Controller) ReadDeviceType(jwt jwt_http_router.Jwt, id string) (dt model.DeviceType, err error, code int) {
	tdt, err, code := this.com.GetTechnicalDeviceType(jwt, id)
	if err != nil {
		return tdt, err, code
	}
	sdt, err, code := this.com.GetSemanticDeviceType(jwt, id)
	if err != nil {
		return tdt, err, code
	}
	tdt.DeviceClass = sdt.DeviceClass
	index := map[string]model.Service{}
	for _, service := range sdt.Services {
		index[service.Id] = service
	}
	for i, service := range tdt.Services {
		service.Functions = index[service.Id].Functions
		service.Aspects = index[service.Id].Aspects
		tdt.Services[i] = service
	}
	return tdt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeCreate(jwt jwt_http_router.Jwt, dt model.DeviceType) (model.DeviceType, error, int) {
	dt.GenerateId()
	err, code := this.com.ValidateDeviceType(jwt, dt)
	if err != nil {
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, jwt.UserId)
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeUpdate(jwt jwt_http_router.Jwt, id string, dt model.DeviceType) (model.DeviceType, error, int) {
	if dt.Id != id {
		return dt, errors.New("device id in body unequal to device id in request endpoint"), http.StatusBadRequest
	}
	err, code := this.com.PermissionCheckForDeviceType(jwt, id, "w")
	if err != nil {
		return dt, err, code
	}
	err, code = this.com.ValidateDeviceType(jwt, dt)
	if err != nil {
		return dt, err, code
	}
	err = this.publisher.PublishDeviceType(dt, jwt.UserId)
	if err != nil {
		return dt, err, http.StatusInternalServerError
	}
	return dt, nil, http.StatusOK
}

func (this *Controller) PublishDeviceTypeDelete(jwt jwt_http_router.Jwt, id string) (error, int) {
	err, code := this.com.PermissionCheckForDeviceType(jwt, id, "w")
	if err != nil {
		return err, code
	}
	err = this.publisher.PublishDeviceDelete(id, jwt.UserId)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
