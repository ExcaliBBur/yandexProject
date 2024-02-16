package entity

type Expression struct {
	Id         int     `json:"id"`
	Result     float64 `json:"result"`
	Expression string  `json:"expression"`
	IsFinished bool    `json:"is_finished"`
	IsError    bool    `json:"is_error"`
}
