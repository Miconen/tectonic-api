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
		Summary:     "Add a new teammate to record by boss",
		Tags:        []string{"Team"},
	}, s.AddTeammateByBoss)

	huma.Register(api, huma.Operation{
		OperationID: "add-teammate-by-record-id",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/teams/id/{record_id}",
		Summary:     "Add a new teammate to record by record ID",
		Tags:        []string{"Team"},
	}, s.AddTeammateByRecordId)

	huma.Register(api, huma.Operation{
		OperationID: "remove-teammate-by-boss",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/teams/boss/{boss}",
		Summary:     "Remove a teammate from record by boss",
		Tags:        []string{"Team"},
	}, s.RemoveTeammateByBoss)

	huma.Register(api, huma.Operation{
		OperationID: "remove-teammate-by-record-id",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/teams/id/{record_id}",
		Summary:     "Remove a teammate from record by record ID",
		Tags:        []string{"Team"},
	}, s.RemoveTeammateByRecordId)
}
