package main

import (
	"fmt"
	"log"
	"net/http"

	"uptime-monitor-backend/cmd/api/factory"
	"uptime-monitor-backend/cmd/api/middleware"
	"uptime-monitor-backend/internal/auth"
	"uptime-monitor-backend/internal/config"
	"uptime-monitor-backend/pkg/database"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
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

	authHandler := factory.NewAuthHandler(db, &cfg)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	jwtAdapter := auth.NewJWTAdapter(cfg.JWTSecret)
	authMiddleware := middleware.NewAuthMiddleware(jwtAdapter)

	r.Use(authMiddleware.Authorize)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Uptime Monitor API."))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
	})

	addr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
