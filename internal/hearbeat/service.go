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

func (s *Service) List24hBymonitor_id(monitor_id string) ([]Heartbeat, error) {
	heartbeats, err := s.repo.List24hBymonitor_id(monitor_id)
	if err != nil {
		return nil, fmt.Errorf("list heartbeats by monitor id: %w", err)
	}
	return heartbeats, nil
}
