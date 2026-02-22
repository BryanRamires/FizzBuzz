package stats

import "errors"

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

func (s *Service) Record(k Key) {
	s.repo.Inc(k)
}

func (s *Service) MostFrequent() (Top, bool) {
	return s.repo.Top()
}
