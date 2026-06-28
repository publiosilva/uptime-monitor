package worker

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"uptime-monitor-backend/internal/hearbeat"
	"uptime-monitor-backend/internal/monitor"
)

type Dispatcher struct {
	prober           *Prober
	monitorService   *monitor.Service
	heartbeatService *hearbeat.Service
}

func NewDispatcher(prober *Prober, monitorService *monitor.Service, heartbeatService *hearbeat.Service) *Dispatcher {
	return &Dispatcher{prober: prober, monitorService: monitorService, heartbeatService: heartbeatService}
}

const (
	monitorChanBufferSize = 1000
)

func (d *Dispatcher) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
				if err := d.monitorService.Update(monitor.Monitor{ID: heartbeat.MonitorID, IsUp: heartbeat.IsUp}); err != nil {
					log.Printf("Error updating monitor isUp: %v\n", err)
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

	<-ctx.Done()
	log.Println("Context cancelled, stopping dispatcher.")

	workerWg.Wait()
	log.Println("All workers stopped.")

	close(heartbeatChan)
	log.Println("Heartbeat channel closed.")

	wg.Wait()
	log.Println("Dispatcher stopped.")
}
