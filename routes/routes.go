package routes

import (
	"tectonic-api/handlers"

	_ "tectonic-api/docs"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func AttachV1Routes(r chi.Router, s *handlers.Server) huma.API {
	config := huma.DefaultConfig("Tectonic API", "0.1")
	config.Info.Description = "Functionality provider for Tectonic guild."
	api := humachi.New(r, config)

	RegisterGuildRoutes(api, s)
	RegisterUserRoutes(api, s)
	RegisterTimeRoutes(api, s)
	RegisterTeamRoutes(api, s)
	RegisterEventRoutes(api, s)
	RegisterPointRoutes(api, s)
	RegisterWomRoutes(api, s)
	RegisterRsnRoutes(api, s)
	RegisterLeaderboardRoutes(api, s)
	RegisterAchievementRoutes(api, s)
	RegisterCombatAchievementRoutes(api, s)
	RegisterMiscRoutes(api, s)

	return api
}
