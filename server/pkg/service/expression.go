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

	var stack utility.Stack
	var taskId = 1

	for _, char := range strings.Split(expression.Expression, " ") {
		if char == "+" || char == "-" || char == "*" || char == "/" {

			op2, _ := stack.Pop()
			op1, _ := stack.Pop()

			if char == "/" && op2 == 0 {
				s.updateExpression(expression, 0, true, true)
				return errors.New("divide by zero")
			}

			expr := fmt.Sprintf("%f%s%f", op1, char, op2)
			var task entity.Task = entity.Task{Operand1: op1, Operand2: op2,
				Expression: expr, Operator: char, TaskId: taskId, ExpressionId: expression.Id}
			s.repo.CreateTask(task)

			if err := s.SendMessage(encode(task)); err != nil {
				log.Fatalf("Can not send to queue expression %s with expressionId %d taskId %d",
					task.Expression, task.ExpressionId, task.TaskId)
				return err
			}
			channel := s.repo.GetChannel()
			res := s.getResult(channel, expression, taskId, task, stopChannel)
			stack.Push(res)
			taskId += 1
		} else {
			num, _ := strconv.ParseFloat(char, 64)
			stack.Push(num)
		}
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
		log.Printf("Waiting for exprID %d with taskID %d", expression.Id, taskId)

		if task.ExpressionId == expression.Id && task.TaskId == taskId {
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
