package handlers

import (
	"net/http"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"
)

// @Summary		Get all supported achievements
// @Description	Get all possible supported achievements from the database
// @Tags			Achievement
// @Produce		json
// @Success		200			{object}	models.ErrorResponse
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/api/v1/achievements [GET]
func GetAchievements(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	achievements, err := queries.GetAchievements(r.Context())
	ei := database.ClassifyError(err)
	if ei != nil {
		// TODO: Apply proper error response
		handleDatabaseError(*ei, jw, models.ERROR_GUILD_NOT_FOUND)
		return
	}

	// Write JSON response
	jw.WriteResponse(achievements)
}
