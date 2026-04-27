package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterEventRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-events",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/events",
		Summary:     "Get the guild's events",
		Tags:        []string{"Event"},
	}, s.GetEvents)

	huma.Register(api, huma.Operation{
		OperationID: "get-detailed-event",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/events/{event_id}",
		Summary:     "Get the guild's event details",
		Tags:        []string{"Event"},
	}, s.GetDetailedEvent)

	huma.Register(api, huma.Operation{
		OperationID: "register-event",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/events",
		Summary:     "Register a guild event",
		Tags:        []string{"Event"},
	}, s.RegisterEvent)

	huma.Register(api, huma.Operation{
		OperationID: "delete-event",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/events/{event_id}",
		Summary:     "Delete a guild event",
		Tags:        []string{"Event"},
	}, s.DeleteEvent)

	huma.Register(api, huma.Operation{
		OperationID: "update-event",
		Method:      http.MethodPut,
		Path:        "/api/v1/guilds/{guild_id}/events/{event_id}",
		Summary:     "Update a guild event",
		Tags:        []string{"Event"},
	}, s.UpdateEvent)

	huma.Register(api, huma.Operation{
		OperationID: "register-legacy-event",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/events/legacy",
		Summary:     "Register a legacy event with Discord user IDs (no WOM)",
		Tags:        []string{"Event"},
	}, s.RegisterLegacyEvent)
}
