package routes

import (
	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func AttachV1Routes(r chi.Router, s *handlers.Server) huma.API {
	config := huma.DefaultConfig("Tectonic API", "1.0.0")
	config.Info.Description = "Functionality provider for Tectonic guild.\n\n" +
		"**Tracing:** send `X-Request-ID` (letters, digits, `.`, `_`, `:` or `-`, up to 64 characters) to tag the API's logs " +
		"with your own correlation ID; one is generated when it's missing, and it's always returned in the `X-Request-ID` " +
		"response header. Send `X-Actor-ID` with the Discord user ID behind the request to follow one member's actions in the logs."
	api := humachi.New(r, config)

	RegisterGuildRoutes(api, s)
	RegisterUserRoutes(api, s)
	RegisterRecordRoutes(api, s)
	RegisterTeamRoutes(api, s)
	RegisterEventRoutes(api, s)
	RegisterPointRoutes(api, s)
	RegisterWomRoutes(api, s)
	RegisterRsnRoutes(api, s)
	RegisterLeaderboardRoutes(api, s)
	RegisterAchievementRoutes(api, s)
	RegisterCombatAchievementRoutes(api, s)
	RegisterGuildRankRoutes(api, s)
	RegisterMiscRoutes(api, s)

	return api
}
