package main

import (
	"log"
	"net/http"
	"tectonic-api/routes"
)

// @title		Tectonic API
// @version		0.1
// @description		Functionality provider for Tectonic guild.
// @host		localhost:8080
// @BasePath		/api
func main() {
	router := routes.NewAPIBuilder().AttachV1Routes()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
