package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGuildRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-guild",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}",
		Summary:     "Get a guild by ID",
		Tags:        []string{"Guild"},
	}, s.GetGuild)

	huma.Register(api, huma.Operation{
		OperationID: "create-guild",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds",
		Summary:     "Create / Initialize a guild",
		Tags:        []string{"Guild"},
	}, s.CreateGuild)

	huma.Register(api, huma.Operation{
		OperationID: "update-guild",
		Method:      http.MethodPut,
		Path:        "/api/v1/guilds/{guild_id}",
		Summary:     "Update a guild",
		Tags:        []string{"Guild"},
	}, s.UpdateGuild)

	huma.Register(api, huma.Operation{
		OperationID: "delete-guild",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}",
		Summary:     "Delete a guild",
		Tags:        []string{"Guild"},
	}, s.DeleteGuild)
}
