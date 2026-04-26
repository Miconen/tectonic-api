package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRecordRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-guild-records",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/records",
		Summary:     "Get all guild records",
		Tags:        []string{"Record"},
	}, s.GetGuildRecords)

	huma.Register(api, huma.Operation{
		OperationID: "create-record",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/records",
		Summary:     "Add a new record to guild",
		Tags:        []string{"Record"},
	}, s.CreateRecord)

	huma.Register(api, huma.Operation{
		OperationID: "remove-record",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/records/id/{record_id}",
		Summary:     "Remove a record from guild",
		Tags:        []string{"Record"},
	}, s.RemoveRecord)

	huma.Register(api, huma.Operation{
		OperationID: "clear-boss-records",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/records/{boss}/clear",
		Summary:     "Clear all records for a boss",
		Tags:        []string{"Record"},
	}, s.ClearBossRecords)

	huma.Register(api, huma.Operation{
		OperationID: "revert-top-record",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/records/{boss}/revert",
		Summary:     "Revert top record (delete #1, next best becomes #1)",
		Tags:        []string{"Record"},
	}, s.RevertTopRecord)
}
