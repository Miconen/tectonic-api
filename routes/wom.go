package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterWomRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "end-competition",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/wom/competition/{competition_id}/cutoff/{cutoff}",
		Summary:     "Handle WOM competition end",
		Tags:        []string{"WOM"},
	}, s.EndCompetition)

	huma.Register(api, huma.Operation{
		OperationID: "competition-winners",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/wom/winners/{competition_id}",
		Summary:     "Get competition winners",
		Tags:        []string{"WOM"},
	}, s.CompetitionWinners)

	huma.Register(api, huma.Operation{
		OperationID: "competition-team-position",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/wom/winners/{competition_id}/team/{team}",
		Summary:     "Get competition winners by team name",
		Tags:        []string{"WOM"},
	}, s.CompetitionTeamPosition)
}
