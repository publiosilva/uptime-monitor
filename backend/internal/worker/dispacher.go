package worker

import (
	"context"
	"log"
	"sync"
	"time"
	"uptime-monitor-backend/internal/hearbeat"
	"uptime-monitor-backend/internal/monitor"
	"uptime-monitor-backend/internal/notification"
)

type Dispatcher struct {
	prober           *Prober
	monitorService   *monitor.Service
	heartbeatService *hearbeat.Service
	topicHub         *notification.TopicHub
}

func NewDispatcher(
	prober *Prober,
	monitorService *monitor.Service,
	heartbeatService *hearbeat.Service,
	topicHub *notification.TopicHub,
) *Dispatcher {
	return &Dispatcher{
		prober:           prober,
		monitorService:   monitorService,
		heartbeatService: heartbeatService,
		topicHub:         topicHub,
	}
}

const (
	monitorChanBufferSize = 1000
)

func (d *Dispatcher) Start(ctx context.Context) {
	monitorChan := make(chan monitor.Monitor, monitorChanBufferSize)
	heartbeatChan := make(chan hearbeat.Heartbeat, monitorChanBufferSize)

	wg := sync.WaitGroup{}
	workerWg := sync.WaitGroup{}

	d.prober.StartProbeQueue(ctx, monitorChan, heartbeatChan, &workerWg)

	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case heartbeat := <-heartbeatChan:
				log.Printf("Received heartbeat: %+v\n", heartbeat)
				_, err := d.heartbeatService.Create(heartbeat)
				if err != nil {
					log.Printf("Error creating heartbeat: %v\n", err)
					continue
				}
				mon, err := d.monitorService.FindByID(heartbeat.MonitorID)
				if err != nil {
					log.Printf("Error getting monitor: %v\n", err)
					continue
				}
				if err := d.monitorService.Update(monitor.Monitor{ID: heartbeat.MonitorID, IsUp: heartbeat.IsUp}); err != nil {
					log.Printf("Error updating monitor isUp: %v\n", err)
				}
				if mon.IsUp != heartbeat.IsUp {
					d.topicHub.Broadcast <- notification.TopicMessage{
						Topic: mon.UserID,
						Payload: map[string]any{
							"monitor_name": mon.Name,
							"monitorId":    mon.ID,
							"is_up":        heartbeat.IsUp,
							"timestamp":    time.Now().Unix(),
						},
					}
				}
			}
		}
	})

	wg.Go(func() {
		defer close(monitorChan)

		monitors, err := d.monitorService.ListAll()
		if err != nil {
			log.Printf("Error listing monitors: %v\n", err)
			return
		}

		scheduleWg := sync.WaitGroup{}

		for _, m := range monitors {
			scheduleWg.Add(1)
			go func(target monitor.Monitor) {
				defer scheduleWg.Done()

				ticker := time.NewTicker(time.Duration(target.Frequency) * time.Second)
				defer ticker.Stop()

				// send first tick immediately
				select {
				case <-ctx.Done():
					return
				case monitorChan <- target:
				default:
					log.Printf("Monitor channel is full, skipping target: %+v\n", target)
				}

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						select {
						case monitorChan <- target:
						default:
							log.Printf("Monitor channel is full, skipping target: %+v\n", target)
						}
					}
				}
			}(m)
		}

		scheduleWg.Wait()
		log.Println("All monitor schedulers stopped.")
	})

	workerWg.Wait()
	log.Println("All workers stopped.")

	close(heartbeatChan)
	log.Println("Heartbeat channel closed.")

	wg.Wait()
	log.Println("Dispatcher stopped.")
}
