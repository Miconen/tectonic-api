package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAchievementRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "give-achievement-by-id",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/achievements/{achievement}",
		Summary:     "Give an achievement to user by ID",
		Tags:        []string{"Achievement"},
	}, s.GiveAchievementById)

	huma.Register(api, huma.Operation{
		OperationID: "give-achievement-by-rsn",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/users/rsn/{rsn}/achievements/{achievement}",
		Summary:     "Give an achievement to user by RSN",
		Tags:        []string{"Achievement"},
	}, s.GiveAchievementByRsn)

	huma.Register(api, huma.Operation{
		OperationID: "remove-achievement-by-id",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/achievements/{achievement}",
		Summary:     "Remove an achievement from user by ID",
		Tags:        []string{"Achievement"},
	}, s.RemoveAchievementById)

	huma.Register(api, huma.Operation{
		OperationID: "remove-achievement-by-rsn",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/rsn/{rsn}/achievements/{achievement}",
		Summary:     "Remove an achievement from user by RSN",
		Tags:        []string{"Achievement"},
	}, s.RemoveAchievementByRsn)
}
