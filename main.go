package main

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/config"
	"tectonic-api/database"
	"tectonic-api/handlers"
	"tectonic-api/logging"
	"tectonic-api/routes"
	"tectonic-api/utils"
)

// @title			Tectonic API
// @version		0.1
// @description	Functionality provider for Tectonic guild.
// @host			localhost:8080
// @BasePath		/api
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading application configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(cfg.DatabaseURL)
	logging.Init(cfg)

	conn, err := database.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = database.RunMigrations(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running migrations: %v\n", err)
		os.Exit(1)
	}
	logging.Get().Info("migrations ran")

	wom := utils.NewWomClient(cfg)
	handlers.InitHandlers(conn, wom)
	router := routes.NewAPIBuilder().AttachV1Routes()
	logging.Get().Info("routes registered")

	logging.Get().Info("server listening to requests", "port", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, router)
	if err != nil {
		logging.Get().Error("Server failed to start", "error", err)
	}
}
