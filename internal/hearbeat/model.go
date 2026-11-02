package hearbeat

import "time"

type Heartbeat struct {
	ID           string
	MonitorID    string
	StatusCode   int
	LatencyMs    int
	IsUp         bool
	ErrorMessage string
	CreatedAt    time.Time
}
