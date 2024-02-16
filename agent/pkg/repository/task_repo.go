package repository

import (
	entity "agent/model/entity"
)

type TaskRepo struct {
	ch <-chan entity.Task
}

func NewTaskRepo(ch <-chan entity.Task) *TaskRepo {
	return &TaskRepo{ch: ch}
}

func (r *TaskRepo) GetChannel() <-chan entity.Task {
	return r.ch
}

func (r *TaskRepo) SendMessage(bytes []byte, queueName string, delay int) error {
	return SendMessage(bytes, queueName, delay)
}
