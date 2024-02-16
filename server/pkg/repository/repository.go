package repository

import (
	"database/sql"
	"server/dto"
	entity "server/model/entity"
)

type Expression interface {
	CreateExpression(expression entity.Expression) (int, error)
	UpdateExpression(task entity.Expression) error

	CreateTask(task entity.Task) error
	UpdateTask(task entity.Task) error

	GetExpressions(pageNumber, pageSize int) ([]dto.ExpressionResponse, error)
	GetExpression(id int) (dto.ExpressionResponse, error)
	GetUnfinishedExpressions() ([]entity.Expression, error)

	GetChannel() chan entity.Task

	SendMessage(bytes []byte) error
}

type Worker interface {
	UpsertWorker(entity.Worker) error
	UpdateWorkers(delay int) ([]entity.Worker, error)

	CleanUpWorkers(delay int) error

	GetHeartBeatChannel() chan entity.Worker
}

type Duration interface {
	GetDuration() (entity.Duration, error)
	UpdateDuration(duration entity.Duration) error

	SendMessage(bytes []byte) error
}

type Idempotency interface {
	IsIdempotencyKeyExists(key string) (bool, error)
	CreateIdempotencyKey(key string, id int) error
	GetExpressionId(key string) (int, error)
}

type Repository struct {
	Expression
	Worker
	Duration
	Idempotency
}

func NewRepository(db *sql.DB, channel chan entity.Task, heartBeatChannel chan entity.Worker) *Repository {
	return &Repository{
		Expression:  NewExpressionRepo(db, channel),
		Worker:      NewWorkerRepo(db, heartBeatChannel),
		Duration:    NewDurationRepo(db),
		Idempotency: NewIdempotencyRepo(db),
	}
}
