package repository

import (
	entity "agent/model/entity"
)

type Task interface {
	GetChannel() <-chan entity.Task
	SendMessage(bytes []byte, queueName string, delay int) error
}

type Duration interface {
	GetChannel() <-chan entity.Duration
}

type Repository struct {
	Task
	Duration
}

func NewRepository(taskChannel <-chan entity.Task, durationChannel <-chan entity.Duration) *Repository {
	return &Repository{
		Task:     NewTaskRepo(taskChannel),
		Duration: NewDurationRepo(durationChannel),
	}
}
