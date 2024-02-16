package service

import (
	entity "agent/model/entity"
	"agent/pkg/repository"

	"github.com/spf13/viper"
)

type Task interface {
	Evaluate(numWorkers int, ch <-chan entity.Duration)
	SendMessage(bytes []byte, queueName string, delay int) error
	StartWorker(id int, jobs <-chan entity.Task)
	EvaluateExpression(task entity.Task, id int) error
}

type Duration interface {
	Listen(ch chan entity.Duration)
}

type Service struct {
	Task
	Duration
}

func NewService(repo *repository.Repository) *Service {
	service := Service{
		Task:     NewTaskService(repo.Task),
		Duration: NewDurationService(repo.Duration),
	}

	durationInternalChannel := make(chan entity.Duration, 15)

	go func(ch chan entity.Duration) {
		service.Duration.Listen(ch)
	}(durationInternalChannel)

	go func(ch <-chan entity.Duration) {
		service.Task.Evaluate(viper.GetInt("max_workers"), ch)
	}(durationInternalChannel)

	return &service
}
