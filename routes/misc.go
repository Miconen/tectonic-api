package routes

import (
	"net/http"

	"tectonic-api/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterMiscRoutes(api huma.API, s *handlers.Server) {
	huma.Register(api, huma.Operation{
		OperationID: "get-bosses",
		Method:      http.MethodGet,
		Path:        "/api/v1/bosses",
		Summary:     "Get all bosses",
		Tags:        []string{"Miscellaneous"},
	}, s.GetBosses)

	huma.Register(api, huma.Operation{
		OperationID: "get-categories",
		Method:      http.MethodGet,
		Path:        "/api/v1/categories",
		Summary:     "Get all categories",
		Tags:        []string{"Miscellaneous"},
	}, s.GetCategories)

	huma.Register(api, huma.Operation{
		OperationID: "get-achievements",
		Method:      http.MethodGet,
		Path:        "/api/v1/achievements",
		Summary:     "Get all supported achievements",
		Tags:        []string{"Achievement"},
	}, s.GetAchievements)
}
