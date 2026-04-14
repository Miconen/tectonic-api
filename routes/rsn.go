package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRsnRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "create-rsn",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/rsns",
		Summary:     "Link an RSN to a user",
		Tags:        []string{"RSN"},
	}, s.CreateRSN)

	huma.Register(api, huma.Operation{
		OperationID: "remove-rsn",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/rsns/{rsn}",
		Summary:     "Remove RSN from guild and user",
		Tags:        []string{"RSN"},
	}, s.RemoveRSN)
}
