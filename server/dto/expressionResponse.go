package dto

import "time"

type ExpressionResponse struct {
	Id         int       `json:"id"`
	Result     float64   `json:"result"`
	Expression string    `json:"expression"`
	IsFinished bool      `json:"is_finished"`
	DateStart  time.Time `json:"date_start"`
	DateFinish time.Time `json:"date_finish"`
	IsError    bool      `json:"is_error"`
}
