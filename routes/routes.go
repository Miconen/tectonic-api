package routes

import (
	"net/http"
	"tectonic-api/handlers"
	"tectonic-api/logging"
	"tectonic-api/middleware"
	"tectonic-api/utils"

	"github.com/gorilla/mux"

	_ "tectonic-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type APIBuilder struct {
	router *mux.Router
	server *handlers.Server
}

func NewAPIBuilder(srv *handlers.Server) *APIBuilder {
	return &APIBuilder{
		router: mux.NewRouter(),
		server: srv,
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
	r.Use(
		logging.LoggingHandler,
		middleware.CORS,
		middleware.RateLimit,
		middleware.Authentication(b.server.Config()),
	)

	// Non-guild functionality
	r.HandleFunc("/bosses", b.server.GetBosses).Methods("GET")
	r.HandleFunc("/categories", b.server.GetCategories).Methods("GET")
	r.HandleFunc("/achievements", b.server.GetAchievements).Methods("GET")

	// Guilds
	guildsRouter := r.PathPrefix("/guilds").Subrouter()
	guildsRouter.HandleFunc("", b.server.CreateGuild).Methods("POST")
	guildsRouter.HandleFunc("/{guild_id}", b.server.UpdateGuild).Methods("PUT")
	guildsRouter.HandleFunc("/{guild_id}", b.server.GetGuild).Methods("GET")
	guildsRouter.HandleFunc("/{guild_id}", b.server.DeleteGuild).Methods("DELETE")

	// Leaderboard
	guildsRouter.HandleFunc("/{guild_id}/leaderboard", b.server.GetLeaderboard).Methods("GET")

	// Events
	eventsRouter := guildsRouter.PathPrefix("/{guild_id}/events").Subrouter()
	eventsRouter.HandleFunc("", b.server.GetEvents).Methods("GET")
	eventsRouter.HandleFunc("", b.server.RegisterEvent).Methods("POST")
	eventsRouter.HandleFunc("/{event_id}", b.server.DeleteEvent).Methods("DELETE")

	// Times
	timesRouter := guildsRouter.PathPrefix("/{guild_id}/times").Subrouter()
	timesRouter.HandleFunc("", b.server.GetGuildTimes).Methods("GET")
	timesRouter.HandleFunc("", b.server.CreateTime).Methods("POST")
	timesRouter.HandleFunc("/{time_id}", b.server.RemoveTime).Methods("DELETE")

	// WOM Events
	womRouter := guildsRouter.PathPrefix("/{guild_id}/wom").Subrouter()
	womRouter.HandleFunc("/competition/{competition_id}/cutoff/{cutoff}", b.server.EndCompetition).Methods("GET")
	womRouter.HandleFunc("/winners/{competition_id}", b.server.CompetitionWinners).Methods("GET")
	womRouter.HandleFunc("/winners/{competition_id}/team/{team}", b.server.CompetitionTeamPosition).Methods("GET")

	// Users
	usersRouter := guildsRouter.PathPrefix("/{guild_id}/users").Subrouter()
	usersRouter.HandleFunc("", b.server.CreateUser).Methods("POST")
	usersRouter.HandleFunc("/{user_id}/times", b.server.GetUserTimes).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/events", b.server.GetUserEvents).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/achievements", b.server.GetUserAchievements).Methods("GET")
	usersRouter.HandleFunc("/rsn/{rsns}", b.server.GetUsersByRsn).Methods("GET")
	usersRouter.HandleFunc("/wom/{wom_ids}", b.server.GetUsersByWom).Methods("GET")
	usersRouter.HandleFunc("/{user_ids}", b.server.GetUsersById).Methods("GET")
	usersRouter.HandleFunc("/rsn/{rsn}", b.server.RemoveUserByRsn).Methods("DELETE")
	usersRouter.HandleFunc("/wom/{wom_id}", b.server.RemoveUserByWom).Methods("DELETE")
	usersRouter.HandleFunc("/{user_id}", b.server.RemoveUserById).Methods("DELETE")

	// Achievements
	usersRouter.HandleFunc("/rsn/{rsn}/achievements/{achievement}", b.server.GiveAchievementByRsn).Methods("POST")
	usersRouter.HandleFunc("/rsn/{rsn}/achievements/{achievement}", b.server.RemoveAchievementByRsn).Methods("DELETE")
	usersRouter.HandleFunc("/{user_id}/achievements/{achievement}", b.server.GiveAchievementById).Methods("POST")
	usersRouter.HandleFunc("/{user_id}/achievements/{achievement}", b.server.RemoveAchievementById).Methods("DELETE")

	// RSN
	rsnsRouter := usersRouter.PathPrefix("/{user_id}/rsns").Subrouter()
	rsnsRouter.HandleFunc("", b.server.CreateRSN).Methods("POST")
	rsnsRouter.HandleFunc("/{rsn}", b.server.RemoveRSN).Methods("DELETE")

	// Points
	pointsRouter := usersRouter.PathPrefix("/{user_ids}/points").Subrouter()
	pointsRouter.HandleFunc("/custom/{points}", b.server.UpdatePointsCustom).Methods("PUT")
	pointsRouter.HandleFunc("/{point_event}", b.server.UpdatePoints).Methods("PUT")

	return b.router
}
