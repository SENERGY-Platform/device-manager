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

package api

import (
	"github.com/SENERGY-Platform/device-manager/lib/api/util"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/service-commons/pkg/accesslog"
	"log"
	"net/http"
	"reflect"
)

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate swag init --instanceName devicemanager -o ../../docs --parseDependency -d . -g api.go

type EndpointMethod = func(config config.Config, router *http.ServeMux, ctrl Controller)

var endpoints = []interface{}{} //list of objects with EndpointMethod

func Start(config config.Config, control Controller) (srv *http.Server, err error) {
	log.Println("start api")
	router := GetRouter(config, control)
	log.Println("listen on port", config.ServerPort)
	srv = &http.Server{Addr: ":" + config.ServerPort, Handler: router}
	go func() { log.Println(srv.ListenAndServe()) }()
	return srv, nil
}

// GetRouter doc
// @title         Device-Manager API
// @version       0.1
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func GetRouter(config config.Config, control Controller) http.Handler {
	handler := GetRouterWithoutMiddleware(config, control)
	handler = util.NewCors(handler)
	handler = accesslog.New(handler)
	if config.EditForward != "" && config.EditForward != "-" {
		handler = util.NewConditionalForward(handler, config.EditForward, func(r *http.Request) bool {
			return r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete
		})
	}
	return handler
}

func GetRouterWithoutMiddleware(config config.Config, command Controller) http.Handler {
	router := http.NewServeMux()
	log.Println("add heart beat endpoint")
	router.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	for _, e := range endpoints {
		for name, call := range getEndpointMethods(e) {
			log.Println("add endpoint " + name)
			call(config, router, command)
		}
	}
	return router
}

func getEndpointMethods(e interface{}) map[string]func(config config.Config, router *http.ServeMux, ctrl Controller) {
	result := map[string]EndpointMethod{}
	objRef := reflect.ValueOf(e)
	methodCount := objRef.NumMethod()
	for i := 0; i < methodCount; i++ {
		m := objRef.Method(i)
		f, ok := m.Interface().(EndpointMethod)
		if ok {
			name := getTypeName(objRef.Type()) + "::" + objRef.Type().Method(i).Name
			result[name] = f
		}
	}
	return result
}

func getTypeName(t reflect.Type) (res string) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
