package main

import (
	"log"
	"uptime-monitor-backend/internal/config"
	"uptime-monitor-backend/internal/hearbeat"
	"uptime-monitor-backend/internal/monitor"
	"uptime-monitor-backend/internal/worker"
	"uptime-monitor-backend/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	prober := worker.NewProber(50)
	monitorService := monitor.NewService(monitor.NewRepository(db))
	heartbeatService := hearbeat.NewService(hearbeat.NewRepository(db))
	dispatcher := worker.NewDispatcher(prober, monitorService, heartbeatService)
	dispatcher.Start()
}
