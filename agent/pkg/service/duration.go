package service

import (
	entity "agent/model/entity"
	"agent/pkg/repository"
)

type DurationService struct {
	repo repository.Duration
}

func NewDurationService(repo repository.Duration) *DurationService {
	return &DurationService{repo: repo}
}

func (s *DurationService) Listen(ch chan entity.Duration) {
	for message := range s.repo.GetChannel() {
		ch <- message
	}
}
