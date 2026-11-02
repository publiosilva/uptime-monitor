package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"uptime-monitor-backend/cmd/worker/middleware"
	"uptime-monitor-backend/internal/auth"
	"uptime-monitor-backend/internal/config"
	"uptime-monitor-backend/internal/hearbeat"
	"uptime-monitor-backend/internal/monitor"
	"uptime-monitor-backend/internal/notification"
	"uptime-monitor-backend/internal/worker"
	"uptime-monitor-backend/pkg/database"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	topicHub := notification.NewTopicHub()
	prober := worker.NewProber(50)
	monitorService := monitor.NewService(monitor.NewRepository(db))
	heartbeatService := hearbeat.NewService(hearbeat.NewRepository(db))
	dispatcher := worker.NewDispatcher(prober, monitorService, heartbeatService, topicHub)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	jwtAdapter := auth.NewJWTAdapter(cfg.JWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(jwtAdapter)

	go dispatcher.Start(ctx)

	go topicHub.Start(ctx)

	r.With(authMiddleware.AuthorizeWebSocket).HandleFunc("/ws/states", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}

		sub := notification.TopicSubscription{Topic: userID, Conn: conn}
		topicHub.Subscribe <- sub

		go func() {
			defer func() {
				topicHub.Unsubscribe <- sub
				conn.Close()
			}()

			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})

	addr := fmt.Sprintf(":%s", cfg.WSPort)

	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		log.Printf("Notification engine listening on %s", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Context cancelled, stopping worker.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown: %v", err)
	}
	log.Println("HTTP server shut down gracefully.")
}
