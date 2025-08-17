package handlers

import (
	"net/http"
	"tectonic-api/database"
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
func (s *Server) GetAchievements(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	achievements, err := s.queries.GetAchievements(r.Context())
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
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
// @Router			/guilds/{guild_id}/users/{user_id}/achievements/{achievement} [POST]
func (s *Server) GiveAchievementById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	err := database.WrapExec(s.queries.GiveAchievementById, r.Context(), database.GiveAchievementByIdParams{
		UserID:          p["user_id"],
		AchievementName: p["achievement"],
		GuildID:         p["guild_id"],
	})
	if err != nil {
		s.handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(http.NoBody)
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
// @Router			/guilds/{guild_id}/users/rsn/{rsn}/achievements/{achievement} [POST]
func (s *Server) GiveAchievementByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	err := database.WrapExec(s.queries.GiveAchievementByRsn, r.Context(), database.GiveAchievementByRsnParams{
		Rsn:             p["rsn"],
		AchievementName: p["achievement"],
		GuildID:         p["guild_id"],
	})
	if err != nil {
		s.handleDatabaseError(*err, jw)
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
// @Router			/guilds/{guild_id}/users/{user_id}/achievements/{achievement} [DELETE]
func (s *Server) RemoveAchievementById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	err := database.WrapExec(s.queries.RemoveAchievementById, r.Context(), database.RemoveAchievementByIdParams{
		UserID:          p["user_id"],
		AchievementName: p["achievement"],
		GuildID:         p["guild_id"],
	})
	if err != nil {
		s.handleDatabaseError(*err, jw)
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
// @Router			/guilds/{guild_id}/users/rsn/{rsn}/achievements/{achievement} [DELETE]
func (s *Server) RemoveAchievementByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	err := database.WrapExec(s.queries.RemoveAchievementByRsn, r.Context(), database.RemoveAchievementByRsnParams{
		Rsn:             p["rsn"],
		AchievementName: p["achievement"],
		GuildID:         p["guild_id"],
	})
	if err != nil {
		s.handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(http.NoBody)
}
