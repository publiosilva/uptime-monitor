package graph

import (
	"time"
	"uptime-monitor-backend/internal/graph/model"
	"uptime-monitor-backend/internal/hearbeat"
)

func calculateUptimePercentage(heartbeats []hearbeat.Heartbeat) float64 {
	if len(heartbeats) == 0 {
		return 0
	}

	uptime := 0
	for _, heartbeat := range heartbeats {
		if heartbeat.IsUp {
			uptime++
		}
	}
	return (float64(uptime) / float64(len(heartbeats))) * 100
}

func calculateAverageLatency(heartbeats []hearbeat.Heartbeat) float64 {
	if len(heartbeats) == 0 {
		return 0
	}

	latency := 0
	for _, heartbeat := range heartbeats {
		latency += heartbeat.LatencyMs
	}
	return float64(latency) / float64(len(heartbeats))
}

func toModelHeartbeats(heartbeats []hearbeat.Heartbeat) []*model.Heartbeat {
	modelHeartbeats := make([]*model.Heartbeat, 0, len(heartbeats))
	for _, heartbeat := range heartbeats {
		modelHeartbeats = append(modelHeartbeats, &model.Heartbeat{
			ID:         heartbeat.ID,
			StatusCode: int32(heartbeat.StatusCode),
			LatencyMs:  int32(heartbeat.LatencyMs),
			IsUp:       heartbeat.IsUp,
			CreatedAt:  heartbeat.CreatedAt.Format(time.RFC3339),
		})
	}
	return modelHeartbeats
}
