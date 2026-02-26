package stats

import (
	"context"
	"errors"
)

// Service holds stats business rules and abstracts the storage backend.
type Service struct {
	repo Repository
}

func NewService(repo Repository) (*Service, error) {
	if repo == nil {
		return nil, errors.New("stats: nil repository")
	}
	return &Service{repo: repo}, nil
}

func (s *Service) Record(ctx context.Context, k Key) error {
	return s.repo.Inc(ctx, k)
}

func (s *Service) MostFrequent(ctx context.Context) (Top, bool, error) {
	return s.repo.Top(ctx)
}
