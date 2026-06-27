package main

import (
	"fmt"
	"log"
	"net/http"

	"uptime-monitor-backend/cmd/api/factory"
	"uptime-monitor-backend/cmd/api/middleware"
	"uptime-monitor-backend/internal/auth"
	"uptime-monitor-backend/internal/config"
	"uptime-monitor-backend/internal/graph"
	"uptime-monitor-backend/internal/hearbeat"
	"uptime-monitor-backend/internal/monitor"
	"uptime-monitor-backend/pkg/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
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

	monitorRepo := monitor.NewRepository(db)
	monitorService := monitor.NewService(monitorRepo)
	hearbeatRepo := hearbeat.NewRepository(db)
	hearbeatService := hearbeat.NewService(hearbeatRepo)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		MonitorService:  monitorService,
		HearbeatService: hearbeatService,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	authHandler := factory.NewAuthHandler(db, &cfg)
	monitorHandler := factory.NewMonitorHandler(db)

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

		r.Post("/monitors", monitorHandler.Create)
		r.Get("/monitors", monitorHandler.List)
		r.Delete("/monitors/{id}", monitorHandler.Delete)
	})

	r.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	addr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("listening on %s", addr)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.APIPort)
	log.Fatal(http.ListenAndServe(addr, r))
}
