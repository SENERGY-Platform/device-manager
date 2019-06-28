package model

type Device struct {
	Id           string `json:"id"`
	LocalId      string `json:"local_id"`
	Name         string `json:"name"`
	DeviceTypeId string `json:"device_type_id"`
}

type DeviceType struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Services    []Service   `json:"services"`
	DeviceClass DeviceClass `json:"device_class"`
}

type Service struct {
	Id          string     `json:"id"`
	LocalId     string     `json:"local_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Aspects     []Aspect   `json:"aspects"`
	ProtocolId  string     `json:"protocol_id"`
	Inputs      []Content  `json:"inputs"`
	Outputs     []Content  `json:"outputs"`
	Functions   []Function `json:"functions"`
}

type VariableType string

const (
	String  VariableType = "http://www.w3.org/2001/XMLSchema#string"
	Integer VariableType = "http://www.w3.org/2001/XMLSchema#integer"
	Float   VariableType = "http://www.w3.org/2001/XMLSchema#decimal"
	Boolean VariableType = "http://www.w3.org/2001/XMLSchema#boolean"

	Array     VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#Array"     //array with predefined length where each element can be of a different type
	Structure VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#structure" //object with predefined fields where each field can be of a different type
	Map       VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#map"       //object/map where each element has to be of the same type but the key can change
	List      VariableType = "http://www.sepl.wifa.uni-leipzig.de/onto/device-repo#list"      //array where each element has to be of the same type but the length can change
)

type Variable struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Type         VariableType `json:"type"`
	SubVariables []Variable   `json:"sub_variables"`
}
