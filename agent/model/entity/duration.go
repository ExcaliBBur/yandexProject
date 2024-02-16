package entity

type Duration struct {
	PlusDuration      int `json:"plus_duration" validate:"min=1,max=10000"`
	MinusDuration     int `json:"minus_duration" validate:"min=1,max=10000"`
	MulDuration       int `json:"mul_duration" validate:"min=1,max=10000"`
	DivDuration       int `json:"div_duration" validate:"min=1,max=10000"`
	HeartBeatDuration int `json:"heartbeat_duration" validate:"min=1,max=100"`
}
