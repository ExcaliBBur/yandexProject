package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"log"
	"server/dto"
	entity "server/model/entity"
	"server/pkg/repository"
	"server/utility"
	"strconv"
	"strings"
	"time"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

type ExpressionService struct {
	repo repository.Expression
}

func NewExpressionService(repo repository.Expression) *ExpressionService {
	return &ExpressionService{repo: repo}
}

func (s *ExpressionService) CreateExpression(expression entity.Expression) (int, error) {
	return s.repo.CreateExpression(expression)
}

func (s *ExpressionService) updateExpression(expression entity.Expression, res float64, is_finished, is_error bool) {
	expression.Result = res
	expression.IsError = is_error
	expression.IsFinished = is_finished
	s.repo.UpdateExpression(expression)
}

func (s *ExpressionService) ParseExpression(expression entity.Expression) (string, error) {
	if strings.Contains(expression.Expression, "**") || strings.Contains(expression.Expression, "^") {
		return "", errors.New("power calculation is prohibited")
	}

	_, err := parser.ParseExpr(expression.Expression)
	if err != nil {
		return "", err
	}

	infixTokens, err := shuntingYard.Scan(expression.Expression)
	if err != nil {
		return "", err
	}
	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		return "", err
	}

	var postfix string
	for _, t := range postfixTokens {
		str := fmt.Sprintf("%v", t.Value)
		postfix += str + " " // Конкатенация - зло.
	}

	return postfix, nil
}

func (s *ExpressionService) EvaluateAndUpdateExpression(expression entity.Expression) error {
	expression.Expression, _ = s.ParseExpression(expression)
	stopChannel := make(chan bool)

	var tasks []entity.Task = splitExpression(expression)
	for _, task := range tasks {
		s.repo.CreateTask(task)
	}
	var stack utility.Stack

	for taskId, task := range tasks {

		if strings.Contains(task.Expression, "$") {
			if strings.Count(task.Expression, "$") == 1 {
				op1, _ := stack.Pop()
				task.Expression = strings.Replace(task.Expression, "$", fmt.Sprintf("%f", op1), 1)
				task.Operand1 = op1
			} else {
				op2, _ := stack.Pop()
				op1, _ := stack.Pop()
				task.Operand1 = op1
				task.Operand2 = op2
				task.Expression = strings.Replace(task.Expression, "$", fmt.Sprintf("%f", op1), 1)
				task.Expression = strings.Replace(task.Expression, "$", fmt.Sprintf("%f", op2), 1)
			}
		}

		if task.Operator == "/" && task.Operand2 == 0 {
			s.updateExpression(expression, 0, true, true)
			return errors.New("divide by zero")
		}

		if err := s.SendMessage(encode(task)); err != nil {
			log.Fatalf("Can not send to queue expression %s with expressionId %d taskId %d",
				task.Expression, task.ExpressionId, task.TaskId)
			return err
		}

		channel := s.repo.GetChannel()
		stack.Push(s.getResult(channel, expression, taskId, task, stopChannel))
	}

	res, _ := stack.Pop()
	s.updateExpression(expression, res, true, false)

	return nil
}

func (s *ExpressionService) getResult(
	channel chan entity.Task, expression entity.Expression,
	taskId int, task entity.Task, stopChannel chan bool) float64 {

	time_start := time.Now()
	go s.checkResultTimeOut(stopChannel, time_start, task)

	for task := range channel {
		log.Printf("Waiting for exprID %d with taskID %d", expression.Id, taskId+1)

		if task.ExpressionId == expression.Id && task.TaskId == taskId+1 {
			log.Printf("Got task %s with result %f exprId %d taskId %d",
				task.Expression, task.Result, task.ExpressionId, task.TaskId)

			go func() {
				stopChannel <- true
			}()
			s.repo.UpdateTask(task)
			return task.Result
		} else {
			if task.ExpressionId != expression.Id {
				channel <- task
				go func() {
					stopChannel <- true
				}()
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
	return 0
}

func (s *ExpressionService) checkResultTimeOut(stopChannel chan bool, time_start time.Time, task entity.Task) {
	for {
		select {
		case <-stopChannel:
			return
		default:
			if time.Since(time_start).Seconds() > 20 {
				log.Printf("TaskId %d result timeout: %fs", task.TaskId, time.Since(time_start).Seconds())
				if err := s.SendMessage(encode(task)); err != nil {
					log.Fatalf("Can not send to queue expression %s with expressionId %d taskId %d",
						task.Expression, task.ExpressionId, task.TaskId)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}
}

func (s *ExpressionService) SendMessage(bytes []byte) error {
	return s.repo.SendMessage(bytes)
}

func (s *ExpressionService) GetExpressions(pageNumber, pageSize int) ([]dto.ExpressionResponse, error) {
	return s.repo.GetExpressions(pageNumber, pageSize)
}

func (s *ExpressionService) GetExpression(id int) (dto.ExpressionResponse, error) {
	return s.repo.GetExpression(id)
}

func (s *ExpressionService) getUnfinishedExpressions() ([]entity.Expression, error) {
	return s.repo.GetUnfinishedExpressions()
}

func (s *ExpressionService) EvaluateUnfinishedExpressions() error {
	expressions, err := s.getUnfinishedExpressions()

	if err != nil {
		return err
	}

	for _, expression := range expressions {
		go func(expr entity.Expression) {
			s.EvaluateAndUpdateExpression(expr)
		}(expression)
	}

	return nil
}

func encode(task entity.Task) []byte {

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(task)

	return reqBodyBytes.Bytes()
}

func splitExpression(expression entity.Expression) []entity.Task {
	var arr []entity.Task
	var stack utility.Stack
	var taskId int = 1
	for _, char := range strings.Split(expression.Expression, " ") {
		if char == "+" || char == "-" || char == "*" || char == "/" {
			var op1 float64
			var op2 float64
			var isOperand1 bool
			var isOperand2 bool

			op2, isOperand2 = stack.Pop()
			op1, isOperand1 = stack.Pop()
			var expr string

			if isOperand1 {
				str := fmt.Sprintf("%f", op1)
				expr += str
			} else {
				expr += "$"
			}

			expr += char

			if isOperand2 {
				str := fmt.Sprintf("%f", op2)
				expr += str
			} else {
				expr += "$"
			}
			var task entity.Task = entity.Task{Operand1: op1, Operand2: op2,
				Expression: expr, Operator: char, TaskId: taskId, ExpressionId: expression.Id}
			taskId += 1
			arr = append(arr, task)
		} else {
			num, _ := strconv.ParseFloat(char, 64)
			stack.Push(num)
		}
	}
	return arr
}
