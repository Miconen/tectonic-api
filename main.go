package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/routes"
)

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

	router := routes.NewAPIBuilder().AttachV1Routes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server will listen on port:", port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
