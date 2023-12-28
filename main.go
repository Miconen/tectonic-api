package main

import (
	"log"
	"net/http"
	"os"
	"tectonic-api/routes"
)

// @title		Tectonic API
// @version		0.1
// @description		Functionality provider for Tectonic guild.
// @host		localhost:8080
// @BasePath		/api
func main() {
	router := routes.NewAPIBuilder().AttachV1Routes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server will listen on port:", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
