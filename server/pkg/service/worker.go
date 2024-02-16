package service

import (
	"log"
	entity "server/model/entity"
	"server/pkg/repository"
	"time"

	"github.com/spf13/viper"
)

type WorkerService struct {
	repo repository.Worker
}

func NewWorkerService(repo repository.Worker) *WorkerService {
	return &WorkerService{repo: repo}
}

func (s *WorkerService) ListenHeartBeats() {
	channel := s.repo.GetHeartBeatChannel()

	for worker := range channel {
		s.UpsertWorker(worker)
	}

}

func (s *WorkerService) UpsertWorker(worker entity.Worker) error {
	return s.repo.UpsertWorker(worker)
}

func (s *WorkerService) UpdateWorkers(delay_seconds int) ([]entity.Worker, error) {
	return s.repo.UpdateWorkers(delay_seconds)
}

func (s *WorkerService) CleanUpWorkers() error {
	delay := viper.GetInt("clean_workers_delay_seconds")
	for {
		err := s.repo.CleanUpWorkers(delay)
		if err != nil {
			log.Fatalf("Can not clean up workers: %s", err.Error())
		}
		time.Sleep(time.Second * time.Duration(delay))
	}
}
