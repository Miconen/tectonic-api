package routes

import (
	"tectonic-api/handlers"
	"tectonic-api/middleware"

	"github.com/gorilla/mux"

	_ "tectonic-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type APIBuilder struct {
	router *mux.Router
}

func NewAPIBuilder() *APIBuilder {
	return &APIBuilder{
		router: mux.NewRouter(),
	}
}

func (b *APIBuilder) AttachV1Routes() *mux.Router {
	r := b.router.PathPrefix("/api/v1").Subrouter()

	r.Use(middleware.Authentication)

	// User
	r.HandleFunc("/user", handlers.GetUser).Methods("GET")
	r.HandleFunc("/user", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/user", handlers.RemoveUser).Methods("DELETE")

	// Users
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	// RSN
	r.HandleFunc("/rsn", handlers.GetRSN).Methods("GET")
	r.HandleFunc("/rsn", handlers.CreateRSN).Methods("POST")
	r.HandleFunc("/rsn", handlers.RemoveRSN).Methods("DELETE")

	// Guild
	r.HandleFunc("/guild", handlers.GetGuild).Methods("GET")
	r.HandleFunc("/guild", handlers.CreateGuild).Methods("POST")
	r.HandleFunc("/guild", handlers.RemoveGuild).Methods("DELETE")

	// Leaderboard
	r.HandleFunc("/leaderboard", handlers.GetLeaderboard).Methods("GET")

	// Time
	r.HandleFunc("/time", handlers.CreateTime).Methods("POST")
	r.HandleFunc("/time", handlers.RemoveTime).Methods("DELETE")

	// Update Times Channel
	r.HandleFunc("/guild/times", handlers.UpdateTimesChannel).Methods("PUT")
	r.HandleFunc("/guild/multiplier", handlers.UpdateMultiplier).Methods("PUT")

	// Serve Swagger UI
	r.PathPrefix("/").Handler(httpSwagger.WrapHandler)

	return r
}
