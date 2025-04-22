package routes

import (
	"net/http"
	"tectonic-api/handlers"
	"tectonic-api/middleware"
	"tectonic-api/utils"

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

	b.router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		jw := utils.NewJsonWriter(w, r, http.StatusOK)
		jw.WriteResponse("pong")
	})

	r := b.router.PathPrefix("/api/v1").Subrouter()
	r.Use(utils.LoggingHandler, middleware.Authentication, handlers.ValidateParameters)

	// Non-guild functionality
	r.HandleFunc("/bosses", handlers.GetBosses).Methods("GET")
	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	r.HandleFunc("/achievements", handlers.GetAchievements).Methods("GET")

	// Achievements
	r.HandleFunc("/achievements/{achievement}/users/{user_id}", handlers.GiveAchievement).Methods("POST")
	r.HandleFunc("/achievements/{achievement}/users/{user_id}", handlers.RemoveAchievement).Methods("DELETE")

	// Guilds
	guildsRouter := r.PathPrefix("/guilds").Subrouter()
	guildsRouter.HandleFunc("", handlers.CreateGuild).Methods("POST")
	guildsRouter.HandleFunc("/{guild_id}", handlers.UpdateGuild).Methods("PUT")
	guildsRouter.HandleFunc("/{guild_id}", handlers.GetGuild).Methods("GET")
	guildsRouter.HandleFunc("/{guild_id}", handlers.DeleteGuild).Methods("DELETE")

	// Leaderboard
	guildsRouter.HandleFunc("/{guild_id}/leaderboard", handlers.GetLeaderboard).Methods("GET")

	// Events
	eventsRouter := guildsRouter.PathPrefix("/{guild_id}/events").Subrouter()
	eventsRouter.HandleFunc("", handlers.GetEvents).Methods("GET")
	eventsRouter.HandleFunc("", handlers.RegisterEvent).Methods("POST")
	eventsRouter.HandleFunc("/{event_id}", handlers.DeleteGuild).Methods("DELETE")

	// Times
	timesRouter := guildsRouter.PathPrefix("/{guild_id}/times").Subrouter()
	timesRouter.HandleFunc("", handlers.GetGuildTimes).Methods("GET")
	timesRouter.HandleFunc("", handlers.CreateTime).Methods("POST")
	timesRouter.HandleFunc("/{time_id}", handlers.RemoveTime).Methods("DELETE")

	// WOM Events
	womRouter := guildsRouter.PathPrefix("/{guild_id}/wom").Subrouter()
	womRouter.HandleFunc("/competition/{competition_id}/cutoff/{cutoff}", handlers.EndCompetition).Methods("GET")
	womRouter.HandleFunc("/winners/{competition_id}", handlers.CompetitionWinners).Methods("GET")
	womRouter.HandleFunc("/winners/{competition_id}/team/{team}", handlers.CompetitionTeamPosition).Methods("GET")

	// Users
	usersRouter := guildsRouter.PathPrefix("/{guild_id}/users").Subrouter()
	usersRouter.HandleFunc("", handlers.CreateUser).Methods("POST")
	usersRouter.HandleFunc("/{user_id}/times", handlers.GetUserTimes).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/events", handlers.GetUserEvents).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/achievements", handlers.GetUserAchievements).Methods("GET")
	usersRouter.HandleFunc("/rsn/{rsns}", handlers.GetUsersByRsn).Methods("GET")
	usersRouter.HandleFunc("/wom/{wom_ids}", handlers.GetUsersByWom).Methods("GET")
	usersRouter.HandleFunc("/{user_ids}", handlers.GetUsersById).Methods("GET")
	usersRouter.HandleFunc("/rsn/{rsn}", handlers.RemoveUserByRsn).Methods("DELETE")
	usersRouter.HandleFunc("/wom/{wom_id}", handlers.RemoveUserByWom).Methods("DELETE")
	usersRouter.HandleFunc("/{user_id}", handlers.RemoveUserById).Methods("DELETE")

	// RSN
	rsnsRouter := usersRouter.PathPrefix("/{user_id}/rsns").Subrouter()
	rsnsRouter.HandleFunc("", handlers.CreateRSN).Methods("POST")
	rsnsRouter.HandleFunc("/{rsn}", handlers.RemoveRSN).Methods("DELETE")

	// Points
	pointsRouter := usersRouter.PathPrefix("/{user_ids}/points").Subrouter()
	pointsRouter.HandleFunc("/custom/{points}", handlers.UpdatePointsCustom).Methods("PUT")
	pointsRouter.HandleFunc("/{point_event}", handlers.UpdatePoints).Methods("PUT")

	return b.router
}
