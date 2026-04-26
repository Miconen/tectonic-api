package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGuildRankRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-guild-ranks",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/ranks",
		Summary:     "Get guild rank tiers",
		Tags:        []string{"Guild Rank"},
	}, s.GetGuildRanks)

	huma.Register(api, huma.Operation{
		OperationID: "create-guild-rank",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/ranks",
		Summary:     "Create a guild rank tier",
		Tags:        []string{"Guild Rank"},
	}, s.CreateGuildRank)

	huma.Register(api, huma.Operation{
		OperationID: "update-guild-rank",
		Method:      http.MethodPut,
		Path:        "/api/v1/guilds/{guild_id}/ranks/{name}",
		Summary:     "Update a guild rank tier",
		Tags:        []string{"Guild Rank"},
	}, s.UpdateGuildRank)

	huma.Register(api, huma.Operation{
		OperationID: "delete-guild-rank",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/ranks/{name}",
		Summary:     "Delete a guild rank tier",
		Tags:        []string{"Guild Rank"},
	}, s.DeleteGuildRank)
}
