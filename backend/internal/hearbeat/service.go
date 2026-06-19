package hearbeat

import "fmt"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(heartbeat Heartbeat) (Heartbeat, error) {
	heartbeat, err := s.repo.Create(heartbeat)
	if err != nil {
		return Heartbeat{}, fmt.Errorf("create heartbeat: %w", err)
	}
	return heartbeat, nil
}
