package model

type DeviceClass struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Function struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Aspect struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
