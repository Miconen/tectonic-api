package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterLeaderboardRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-leaderboard",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/leaderboard",
		Summary:     "Get a guild's leaderboard by ID",
		Tags:        []string{"Leaderboard"},
	}, s.GetLeaderboard)
}
