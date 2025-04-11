package handlers

import (
	"net/http"
	"tectonic-api/database"
	"tectonic-api/utils"
)

//	@Summary		Get all bosses
//	@Description	Get all bosses tracked by the application
//	@Tags			Miscellaneous
//	@Produce		json
//	@Success		200	{object}	models.Guild
//	@Failure		429	{object}	models.Empty
//	@Failure		500	{object}	models.Empty
//	@Router			/api/v1/bosses [GET]
func GetBosses(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	bosses, err := queries.GetBosses(r.Context())
	ei := database.ClassifyError(err)
	if err != nil {
		handleDatabaseError(*ei, jw)
		return
	}

	// Write JSON response
	jw.WriteResponse(bosses)
}

//	@Summary		Get all categories
//	@Description	Get all categories tracked by the application
//	@Tags			Miscellaneous
//	@Produce		json
//	@Success		200	{object}	models.Guild
//	@Failure		429	{object}	models.Empty
//	@Failure		500	{object}	models.Empty
//	@Router			/api/v1/categories [GET]
func GetCategories(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	categories, err := queries.GetCategories(r.Context())
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw)
		return
	}

	// Write JSON response
	jw.WriteResponse(categories)
}
