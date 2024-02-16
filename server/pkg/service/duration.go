package service

import (
	entity "server/model/entity"
	"server/pkg/repository"
)

type DurationService struct {
	repo repository.Duration
}

var (
	delay_seconds int = 5
)

func NewDurationService(repo repository.Duration) *DurationService {
	return &DurationService{repo: repo}
}

func (s *DurationService) SendMessage(bytes []byte) error {
	return s.repo.SendMessage(bytes)
}

func (s *DurationService) GetDelay() int {
	return delay_seconds
}

func (s *DurationService) SetDelay(delay int) {
	delay_seconds = delay
}

func (s *DurationService) UpdateDuration(duration entity.Duration) error {
	s.SetDelay(duration.HeartBeatDuration)
	return s.repo.UpdateDuration(duration)
}

func (s *DurationService) GetDuration() (entity.Duration, error) {
	return s.repo.GetDuration()
}
