package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterPointRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "update-points",
		Method:      http.MethodPut,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_ids}/points/{point_event}",
		Summary:     "Update user(s) points by event",
		Tags:        []string{"Points"},
	}, s.UpdatePoints)

	huma.Register(api, huma.Operation{
		OperationID: "update-points-custom",
		Method:      http.MethodPut,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_ids}/points/custom/{points}",
		Summary:     "Update user(s) points with custom amount",
		Tags:        []string{"Points"},
	}, s.UpdatePointsCustom)

	huma.Register(api, huma.Operation{
		OperationID: "get-point-sources",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/points",
		Summary:     "Get guild point sources",
		Tags:        []string{"Points"},
	}, s.GetPointSources)

	huma.Register(api, huma.Operation{
		OperationID: "update-guild-point-source",
		Method:      http.MethodPut,
		Path:        "/api/v1/guilds/{guild_id}/points/{point_source}/{points}",
		Summary:     "Update a guild point source",
		Tags:        []string{"Points"},
	}, s.UpdateGuildPointSource)
}
