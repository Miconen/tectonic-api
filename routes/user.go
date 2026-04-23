package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterUserRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/api/v1/guilds/{guild_id}/users",
		Summary:     "Create / Initialize a new user",
		Tags:        []string{"User"},
	}, s.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "get-users-by-id",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_ids}",
		Summary:     "Get one or more users by ID(s)",
		Tags:        []string{"User"},
	}, s.GetUsersById)

	huma.Register(api, huma.Operation{
		OperationID: "get-users-by-rsn",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/users/rsn/{rsns}",
		Summary:     "Get one or more users by RSN(s)",
		Tags:        []string{"User"},
	}, s.GetUsersByRsn)

	huma.Register(api, huma.Operation{
		OperationID: "get-users-by-wom",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/users/wom/{wom_ids}",
		Summary:     "Get one or more users by WOM ID(s)",
		Tags:        []string{"User"},
	}, s.GetUsersByWom)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-times",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/times",
		Summary:     "Get user times",
		Tags:        []string{"User"},
	}, s.GetUserTimes)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-events",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/events",
		Summary:     "Get user events",
		Tags:        []string{"User"},
	}, s.GetUserEvents)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-achievements",
		Method:      http.MethodGet,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}/achievements",
		Summary:     "Get user achievements",
		Tags:        []string{"User"},
	}, s.GetUserAchievements)

	huma.Register(api, huma.Operation{
		OperationID: "remove-user-by-id",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/{user_id}",
		Summary:     "Delete a user by User ID",
		Tags:        []string{"User"},
	}, s.RemoveUserById)

	huma.Register(api, huma.Operation{
		OperationID: "remove-user-by-rsn",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/rsn/{rsn}",
		Summary:     "Delete a user by RSN",
		Tags:        []string{"User"},
	}, s.RemoveUserByRsn)

	huma.Register(api, huma.Operation{
		OperationID: "remove-user-by-wom",
		Method:      http.MethodDelete,
		Path:        "/api/v1/guilds/{guild_id}/users/wom/{wom_id}",
		Summary:     "Delete a user by WOM ID",
		Tags:        []string{"User"},
	}, s.RemoveUserByWom)
}
