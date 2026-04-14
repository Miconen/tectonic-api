package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterCombatAchievementRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-guild-combat-achievements",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/combat-achievements",
		Summary:     "Get guild combat achievements",
		Tags:        []string{"CombatAchievement"},
	}, s.GetGuildCombatAchievements)

	huma.Register(api, huma.Operation{
		OperationID: "create-combat-achievement",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/combat-achievements",
		Summary:     "Create a combat achievement for a guild",
		Tags:        []string{"CombatAchievement"},
	}, s.CreateCombatAchievement)

	huma.Register(api, huma.Operation{
		OperationID: "delete-combat-achievement",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/combat-achievements/{ca_name}",
		Summary:     "Delete a combat achievement from a guild",
		Tags:        []string{"CombatAchievement"},
	}, s.DeleteCombatAchievement)

	huma.Register(api, huma.Operation{
		OperationID: "complete-combat-achievement",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/combat-achievements/{ca_name}/complete",
		Summary:     "Complete a combat achievement",
		Tags:        []string{"CombatAchievement"},
	}, s.CompleteCombatAchievement)

	huma.Register(api, huma.Operation{
		OperationID: "give-user-combat-achievement",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/combat-achievements/{ca_name}",
		Summary:     "Grant a combat achievement to a user",
		Tags:        []string{"CombatAchievement"},
	}, s.GiveUserCombatAchievement)

	huma.Register(api, huma.Operation{
		OperationID: "remove-user-combat-achievement",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/combat-achievements/{ca_name}",
		Summary:     "Remove a combat achievement from a user",
		Tags:        []string{"CombatAchievement"},
	}, s.RemoveUserCombatAchievement)
}
