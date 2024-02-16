package service

import (
	"server/pkg/repository"
)

type IdempotencyService struct {
	repo repository.Idempotency
}

func NewIdempotencyService(repo repository.Idempotency) *IdempotencyService {
	return &IdempotencyService{repo: repo}
}

func (s *IdempotencyService) IsIdempotencyKeyExists(key string) (bool, error) {
	return s.repo.IsIdempotencyKeyExists(key)
}

func (s *IdempotencyService) CreateIdempotencyKey(key string, expression_id int) error {
	return s.repo.CreateIdempotencyKey(key, expression_id)
}

func (s *IdempotencyService) GetExpressionId(key string) (int, error) {
	return s.repo.GetExpressionId(key)
}
