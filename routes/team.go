package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterTeamRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "add-teammate-by-boss",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/teams/boss/{boss}",
		Summary:     "Add a new teammate to time by boss",
		Tags:        []string{"Team"},
	}, s.AddTeammateByBoss)

	huma.Register(api, huma.Operation{
		OperationID: "add-teammate-by-run-id",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/teams/id/{run_id}",
		Summary:     "Add a new teammate to time by run ID",
		Tags:        []string{"Team"},
	}, s.AddTeammateByRunId)

	huma.Register(api, huma.Operation{
		OperationID: "remove-teammate-by-boss",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/teams/boss/{boss}",
		Summary:     "Remove a teammate from time by boss",
		Tags:        []string{"Team"},
	}, s.RemoveTeammateByBoss)

	huma.Register(api, huma.Operation{
		OperationID: "remove-teammate-by-run-id",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/teams/id/{run_id}",
		Summary:     "Remove a teammate from time by run ID",
		Tags:        []string{"Team"},
	}, s.RemoveTeammateByRunId)
}
