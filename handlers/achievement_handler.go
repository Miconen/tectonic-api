package handlers

import (
	"net/http"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary		Get all supported achievements
// @Description	Get all possible supported achievements from the database
// @Tags			Achievement
// @Produce		json
// @Success		200			{object}	[]database.Achievement
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
		handleDatabaseError(*ei, jw, models.ERROR_ACHIEVEMENT_NOT_FOUND)
		return
	}

	// Write JSON response
	jw.WriteResponse(achievements)
}

// @Summary		Give an achievement to the user
// @Description	Give an achievement to the user.
// @Tags			Achievement
// @Produce		json
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/achievements/{achievement}/users/{user_id} [POST]
func GiveAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	err := database.WrapExec(queries.GiveAchievement, r.Context(), database.GiveAchievementParams{
		UserID:          p["user_id"],
		AchievementName: p["achievement"],
	})
	if err != nil {
		handleDatabaseError(*err, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Remove an achievement from the user
// @Description	Remove an achievement from the user.
// @Tags			Achievement
// @Produce		json
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/achievements/{achievement}/users/{user_id} [DELETE]
func RemoveAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	err := database.WrapExec(queries.RemoveAchievement, r.Context(), database.RemoveAchievementParams{
		UserID:          p["user_id"],
		AchievementName: p["achievement"],
	})
	if err != nil {
		handleDatabaseError(*err, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}
