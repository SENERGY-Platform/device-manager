package model

type Hub struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Hash           string   `json:"hash"`
	DeviceLocalIds []string `json:"device_local_ids"`
}

type Protocol struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Handler          string            `json:"handler"`
	ProtocolSegments []ProtocolSegment `json:"protocol_segments"`
}

type ProtocolSegment struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Content struct {
	Id                   string                `json:"id"`
	Variable             Variable              `json:"variable"`
	SerializationId      string                `json:"serialization_id"`
	SerializationOptions []SerializationOption `json:"serialization_options"`
	ProtocolSegmentId    string                `json:"protocol_segment_id"`
}

type SerializationOption struct {
	Id         string `json:"id"`
	Option     string `json:"option"`
	VariableId string `json:"variable_id"`
}

type Serialization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
