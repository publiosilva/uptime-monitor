package monitor

import (
	"errors"
	"fmt"
	"time"
)

type ErrRequiredFieldMissing struct {
	Field string
}

func (e ErrRequiredFieldMissing) Error() string {
	return fmt.Sprintf("required field %s is missing", e.Field)
}

var (
	ErrInvalidName      = errors.New("name must be between 1 and 100 characters")
	ErrInvalidURL       = errors.New("url must be a valid http or https URL")
	ErrInvalidMethod    = errors.New("method must be a supported HTTP verb")
	ErrInvalidTimeout   = errors.New("timeout must be between 1 and 30 seconds")
	ErrInvalidFrequency = errors.New("frequency must be between 30 and 86400 seconds")
	ErrMonitorNotFound  = errors.New("monitor not found")
)

type Monitor struct {
	ID        string
	UserID    string
	Name      string
	URL       string
	Method    string
	Timeout   int
	Frequency int
	IsActive  bool
	CreatedAt time.Time
}

type CreateMonitorRequest struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Method    string `json:"method,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
	Frequency int    `json:"frequency,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

type MonitorResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Method    string    `json:"method"`
	Timeout   int       `json:"timeout"`
	Frequency int       `json:"frequency"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (m Monitor) ToResponse() MonitorResponse {
	return MonitorResponse{
		ID:        m.ID,
		Name:      m.Name,
		URL:       m.URL,
		Method:    m.Method,
		Timeout:   m.Timeout,
		Frequency: m.Frequency,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
	}
}
