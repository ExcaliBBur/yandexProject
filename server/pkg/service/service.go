package service

import (
	"log"
	"server/dto"
	entity "server/model/entity"
	"server/pkg/repository"
)

type Expression interface {
	ParseExpression(expression entity.Expression) (string, error)
	CreateExpression(expression entity.Expression) (int, error)
	EvaluateAndUpdateExpression(expression entity.Expression) error

	EvaluateUnfinishedExpressions() error

	GetExpressions(pageNumber, pageSize int) ([]dto.ExpressionResponse, error)
	GetExpression(id int) (dto.ExpressionResponse, error)

	SendMessage(bytes []byte) error
}

type Worker interface {
	UpdateWorkers(delay_milliseconds int) ([]entity.Worker, error)
	UpsertWorker(entity.Worker) error

	CleanUpWorkers() error
	ListenHeartBeats()
}

type Duration interface {
	GetDuration() (entity.Duration, error)
	UpdateDuration(duration entity.Duration) error

	SetDelay(delay int)
	GetDelay() int

	SendMessage(bytes []byte) error
}

type Idempotency interface {
	IsIdempotencyKeyExists(key string) (bool, error)
	CreateIdempotencyKey(key string, expression_id int) error
	GetExpressionId(key string) (int, error)
}

type Service struct {
	Expression
	Worker
	Duration
	Idempotency
}

func NewService(repo *repository.Repository) *Service {
	service := Service{
		Expression:  NewExpressionService(repo.Expression),
		Worker:      NewWorkerService(repo.Worker),
		Duration:    NewDurationService(repo.Duration),
		Idempotency: NewIdempotencyService(repo.Idempotency),
	}

	go func() {
		service.Worker.ListenHeartBeats()
	}()

	go func() {
		service.Worker.CleanUpWorkers()
	}()

	err := service.Expression.EvaluateUnfinishedExpressions()

	if err != nil {
		log.Fatalf("Can not evaluate unfinished expressions: %s", err.Error())
	}
	return &service
}
