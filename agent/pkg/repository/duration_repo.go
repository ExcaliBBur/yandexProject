package repository

import (
	entity "agent/model/entity"
)

type DurationRepo struct {
	ch <-chan entity.Duration
}

func NewDurationRepo(ch <-chan entity.Duration) *DurationRepo {
	return &DurationRepo{ch: ch}
}

func (r *DurationRepo) GetChannel() <-chan entity.Duration {
	return r.ch
}
