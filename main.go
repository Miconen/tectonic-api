package main

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/handlers"
	"tectonic-api/routes"
	"tectonic-api/utils"
)

var log = utils.NewLogger()

// @title			Tectonic API
// @version		0.1
// @description	Functionality provider for Tectonic guild.
// @host			localhost:8080
// @BasePath		/api
func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading application configuration: %v\n", err)
		os.Exit(1)
	}

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
	log.Info("migrations ran")

	wom := utils.NewWomClient(cfg)
	handlers.InitHandlers(conn, wom)
	router := routes.NewAPIBuilder().AttachV1Routes()
	log.Info("routes registered")

	log.Info("server listening to requests", "port", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, router)
	if err != nil {
		log.Error("Server failed to start", "error", err)
	}
}
