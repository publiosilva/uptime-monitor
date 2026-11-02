package worker

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
	"uptime-monitor-backend/internal/hearbeat"
	"uptime-monitor-backend/internal/monitor"
)

type Prober struct {
	client  *http.Client
	workers int
}

func NewProber(workers int) *Prober {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Prober{
		client: &http.Client{
			Transport: transport,
		},
		workers: workers,
	}
}

func (p *Prober) StartProbeQueue(ctx context.Context, monitorChan <-chan monitor.Monitor, heartbeatChan chan<- hearbeat.Heartbeat, wg *sync.WaitGroup) {
	for workerID := range p.workers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case monitor, ok := <-monitorChan:
					if !ok {
						return
					}
					heartbeatChan <- p.probe(ctx, monitor)
				}
			}
		}(workerID)
	}
}

func (p *Prober) probe(ctx context.Context, monitor monitor.Monitor) hearbeat.Heartbeat {
	heartbeat := hearbeat.Heartbeat{
		MonitorID: monitor.ID,
		CreatedAt: time.Now(),
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Duration(monitor.Timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, monitor.Method, monitor.URL, nil)
	if err != nil {
		heartbeat.ErrorMessage = err.Error()
		return heartbeat
	}

	start := time.Now()
	resp, err := p.client.Do(req)
	heartbeat.LatencyMs = int(time.Since(start).Milliseconds())

	if err != nil {
		heartbeat.ErrorMessage = err.Error()
		return heartbeat
	}

	// Ensure body is read and closed to allow connection reuse
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	heartbeat.StatusCode = resp.StatusCode
	heartbeat.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300

	return heartbeat
}
