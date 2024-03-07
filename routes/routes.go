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
	// Serve Swagger UI
	b.router.PathPrefix("/swagger/v1").Handler(httpSwagger.WrapHandler)

	r := b.router.PathPrefix("/api/v1").Subrouter()
	r.Use(middleware.Authentication)

	// Guilds
	guildsRouter := r.PathPrefix("/guilds").Subrouter()
	guildsRouter.HandleFunc("", handlers.CreateGuild).Methods("POST")
	guildsRouter.HandleFunc("/{guild_id}", handlers.UpdateGuild).Methods("PUT")
	guildsRouter.HandleFunc("/{guild_id}", handlers.GetGuild).Methods("GET")
	guildsRouter.HandleFunc("/{guild_id}", handlers.RemoveGuild).Methods("DELETE")

	// Leaderboard
	guildsRouter.HandleFunc("/{guild_id}/leaderboard", handlers.GetLeaderboard).Methods("GET")

	// Times
	timesRouter := guildsRouter.PathPrefix("/{guild_id}/times").Subrouter()
	timesRouter.HandleFunc("", handlers.CreateTime).Methods("POST")
	timesRouter.HandleFunc("/{time_id}", handlers.RemoveTime).Methods("DELETE")

	// Users
	usersRouter := guildsRouter.PathPrefix("/{guild_id}/users").Subrouter()
	usersRouter.HandleFunc("", handlers.GetUsers).Methods("GET")
	usersRouter.HandleFunc("", handlers.CreateUser).Methods("POST")
	usersRouter.HandleFunc("/{user_id}", handlers.GetUser).Methods("GET")
	usersRouter.HandleFunc("/{user_id}", handlers.UpdateUser).Methods("PUT")
	usersRouter.HandleFunc("/{user_id}", handlers.RemoveUser).Methods("DELETE")

	// RSN
	rsnsRouter := usersRouter.PathPrefix("/{user_id}/rsns").Subrouter()
	rsnsRouter.HandleFunc("", handlers.GetRSNs).Methods("GET")
	rsnsRouter.HandleFunc("", handlers.CreateRSN).Methods("POST")
	rsnsRouter.HandleFunc("/{rsn}", handlers.RemoveRSN).Methods("DELETE")

	return b.router
}
