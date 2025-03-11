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
// @version			0.1
// @description		Functionality provider for Tectonic guild.
// @host			localhost:8080
// @BasePath		/api
func main() {
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

	handlers.InitHandlers(conn)
	router := routes.NewAPIBuilder().AttachV1Routes()
	log.Info("routes registered")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("server listening to requests", "port", port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Error("Server failed to start", "error", err)
	}
}
