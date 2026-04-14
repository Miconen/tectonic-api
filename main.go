package main

import (
	"fmt"
	"net/http"
	"os"

	"tectonic-api/config"
	"tectonic-api/database"
	"tectonic-api/handlers"
	"tectonic-api/logging"
	"tectonic-api/middleware"
	"tectonic-api/routes"
	"tectonic-api/utils"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading application configuration: %v\n", err)
		os.Exit(1)
	}

	logging.Init(cfg)

	conn, err := database.InitDB(cfg)
	if err != nil {
		logging.Get().Error("error initializing database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = database.RunMigrations(conn)
	if err != nil {
		logging.Get().Error("error running migrations", "error", err)
		os.Exit(1)
	}

	logging.Get().Info("migrations ran")

	wom := utils.NewWomClient(cfg)

	srv, err := handlers.NewServer(conn, wom, cfg)
	if err != nil {
		logging.Get().Error("error creating server", "error", err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Use(
		logging.LoggingHandler,
		middleware.CORS,
		middleware.RateLimit,
		middleware.Authentication(cfg),
	)

	routes.AttachV1Routes(r, srv)
	logging.Get().Info("routes registered")

	logging.Get().Info("server listening to requests", "port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		logging.Get().Error("Server failed to start", "error", err)
	}
}
