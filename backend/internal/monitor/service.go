package monitor

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

var allowedMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(userID string, req CreateMonitorRequest) (Monitor, error) {
	name := strings.TrimSpace(req.Name)
	monitorURL := strings.TrimSpace(req.URL)
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = "GET"
	}

	timeout := req.Timeout
	if timeout == 0 {
		timeout = 5
	}

	frequency := req.Frequency
	if frequency == 0 {
		frequency = 60
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if err := validateCreateInput(name, monitorURL, method, timeout, frequency); err != nil {
		return Monitor{}, err
	}

	return s.repo.Create(Monitor{
		UserID:    userID,
		Name:      name,
		URL:       monitorURL,
		Method:    method,
		Timeout:   timeout,
		Frequency: frequency,
		IsActive:  isActive,
	})
}

func (s *Service) List(userID string) ([]Monitor, error) {
	monitors, err := s.repo.ListByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("list monitors: %w", err)
	}
	return monitors, nil
}

func (s *Service) ListAll() ([]Monitor, error) {
	monitors, err := s.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("list all monitors: %w", err)
	}
	return monitors, nil
}

func (s *Service) Delete(userID, monitor_id string) error {
	if err := s.repo.DeleteByID(userID, monitor_id); err != nil {
		return err
	}
	return nil
}

func (s *Service) FindByUserAndID(userID, monitor_id string) (Monitor, error) {
	monitor, err := s.repo.FindByUserAndID(userID, monitor_id)
	if err != nil {
		return Monitor{}, fmt.Errorf("find monitor by user and id: %w", err)
	}
	return monitor, nil
}

func (s *Service) FindByID(monitor_id string) (Monitor, error) {
	monitor, err := s.repo.FindByID(monitor_id)
	if err != nil {
		return Monitor{}, fmt.Errorf("find monitor by id: %w", err)
	}
	return monitor, nil
}

func (s *Service) Update(monitor Monitor) error {
	if err := s.repo.Update(monitor); err != nil {
		return fmt.Errorf("update monitor: %w", err)
	}
	return nil
}

func validateCreateInput(name, monitorURL, method string, timeout, frequency int) error {
	if name == "" {
		return ErrRequiredFieldMissing{Field: "name"}
	}
	if len(name) > 100 {
		return ErrInvalidName
	}
	if monitorURL == "" {
		return ErrRequiredFieldMissing{Field: "url"}
	}

	parsedURL, err := url.Parse(monitorURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrInvalidURL
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidURL
	}

	if !slices.Contains(allowedMethods, method) {
		return ErrInvalidMethod
	}

	if timeout < 1 || timeout > 30 {
		return ErrInvalidTimeout
	}

	if frequency < 30 || frequency > 86400 {
		return ErrInvalidFrequency
	}

	return nil
}
