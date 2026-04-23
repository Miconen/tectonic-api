package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterTimeRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-guild-times",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/times",
		Summary:     "Get all guild times",
		Tags:        []string{"Time"},
	}, s.GetGuildTimes)

	huma.Register(api, huma.Operation{
		OperationID: "create-time",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/times",
		Summary:     "Add a new best time to guild",
		Tags:        []string{"Time"},
	}, s.CreateTime)

	huma.Register(api, huma.Operation{
		OperationID: "remove-time",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/times/id/{time_id}",
		Summary:     "Remove time from guild's best times",
		Tags:        []string{"Time"},
	}, s.RemoveTime)

	huma.Register(api, huma.Operation{
		OperationID: "clear-clan-pb",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/times/{boss}/clear",
		Summary:     "Clear best PB time",
		Tags:        []string{"Time"},
	}, s.ClearClanPb)

	huma.Register(api, huma.Operation{
		OperationID: "revert-clan-pb",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/times/{boss}/revert",
		Summary:     "Revert best PB time to second last",
		Tags:        []string{"Time"},
	}, s.RevertClanPb)
}
