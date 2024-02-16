package service

import (
	entity "agent/model/entity"
	"agent/pkg/repository"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"time"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}

var (
	PlusDuration      int = 200
	MinusDuration     int = 200
	MulDuration       int = 200
	DivDuration       int = 200
	HeartBeatDuration int = 5
)

func HandleDuration(ch <-chan entity.Duration) {
	for duration := range ch {
		PlusDuration = duration.PlusDuration
		MinusDuration = duration.MinusDuration
		MulDuration = duration.MulDuration
		DivDuration = duration.DivDuration
		HeartBeatDuration = duration.HeartBeatDuration
		log.Printf("Set new duration: `+`: %dms, `-`: %dms, `*`:%dms, `/`:%dms, HeartBeat:%ds",
			PlusDuration, MinusDuration, MulDuration, DivDuration, HeartBeatDuration)
	}
}

func (s *TaskService) Evaluate(numWorkers int, ch <-chan entity.Duration) {
	hostname, _ := os.Hostname()
	log.Printf("Started server with name %s", hostname)

	go func(channel <-chan entity.Duration) {
		HandleDuration(channel)
	}(ch)

	go func() {
		for {
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(entity.Worker{Hostname: hostname, IsAlive: true, Id: 0})
			s.SendMessage(reqBodyBytes.Bytes(), "SendHeartBeat", HeartBeatDuration)
		}
	}()

	jobs := make(chan entity.Task, 15)

	for w := 1; w <= numWorkers; w++ {
		go s.StartWorker(w, jobs)
	}

	for message := range s.repo.GetChannel() {
		jobs <- message
	}

}

func (s *TaskService) SendMessage(bytes []byte, queueName string, delay int) error {
	return s.repo.SendMessage(bytes, queueName, delay)
}

func (s *TaskService) StartWorker(id int, jobs <-chan entity.Task) {
	for task := range jobs {
		s.EvaluateExpression(task, id)
	}
}

func (s *TaskService) EvaluateExpression(task entity.Task, id int) error {
	log.Printf("Evaluating task %s worker %d", task.Expression, id)

	var res float64

	switch task.Operator {
	case "+":
		res = task.Operand1 + task.Operand2
		time.Sleep(time.Millisecond * time.Duration(PlusDuration))
	case "-":
		res = task.Operand1 - task.Operand2
		time.Sleep(time.Millisecond * time.Duration(MinusDuration))
	case "*":
		res = task.Operand1 * task.Operand2
		time.Sleep(time.Millisecond * time.Duration(MulDuration))
	case "/":
		res = task.Operand1 / task.Operand2
		time.Sleep(time.Millisecond * time.Duration(DivDuration))
	}
	task.Result = res

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(task)

	if err := s.SendMessage(reqBodyBytes.Bytes(), "SendResult", 0); err != nil {
		return err
	}
	log.Printf("Task %s with result %f exprId %d taskId %d in query",
		task.Expression, task.Result, task.ExpressionId, task.TaskId)

	return nil
}
