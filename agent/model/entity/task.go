package entity

type Task struct {
	Expression   string
	Operand1     float64
	Operand2     float64
	Operator     string
	Result       float64
	ExpressionId int
	TaskId       int
}
