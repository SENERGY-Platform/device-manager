package mock

import (
	"encoding/json"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/http/httptest"
)

type PermSearch struct {
	ts *httptest.Server
}

func NewPermSearch() *PermSearch {
	repo := &PermSearch{}

	router := jwt_http_router.New(jwt_http_router.JwtConfig{ForceAuth: true, ForceUser: true})

	router.GET("/jwt/check/:resource/:id/:permission/bool", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		json.NewEncoder(writer).Encode(true)
	})

	repo.ts = httptest.NewServer(router)

	return repo
}

func (this *PermSearch) Stop() {
	this.ts.Close()
}

func (this *PermSearch) Url() string {
	return this.ts.URL
}

