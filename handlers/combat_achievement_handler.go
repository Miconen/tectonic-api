package handlers

import (
	"net/http"

	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary		Get guild combat achievements
// @Description	Get all combat achievements configured for a guild
// @Tags			CombatAchievement
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		200			{object}	[]database.GetGuildCombatAchievementsRow
// @Failure		401			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/combat-achievements [GET]
func (s *Server) GetGuildCombatAchievements(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	cas, err := database.WrapQuery(s.queries.GetGuildCombatAchievements, r.Context(), p["guild_id"])
	if err != nil {
		s.handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(cas)
}

// @Summary		Create a combat achievement for a guild
// @Description	Create a new combat achievement linked to an existing point source
// @Tags			CombatAchievement
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string								true	"Guild ID"
// @Param			body		body		models.CreateCombatAchievementBody	true	"Combat Achievement"
// @Success		201			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		409			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/combat-achievements [POST]
func (s *Server) CreateCombatAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)
	p := mux.Vars(r)

	var body models.CreateCombatAchievementBody
	if err := utils.ParseAndValidateRequestBody(w, r, &body); err != nil {
		return
	}

	err := s.queries.CreateCombatAchievement(r.Context(), database.CreateCombatAchievementParams{
		Name:        body.Name,
		GuildID:     p["guild_id"],
		PointSource: body.PointSource,
	})
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Delete a combat achievement from a guild
// @Description	Delete a combat achievement and all user completions (cascading)
// @Tags			CombatAchievement
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			ca_name		path		string	true	"Combat Achievement Name"
// @Success		204			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/combat-achievements/{ca_name} [DELETE]
func (s *Server) DeleteCombatAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	rows, err := s.queries.DeleteCombatAchievement(r.Context(), database.DeleteCombatAchievementParams{
		Name:    p["ca_name"],
		GuildID: p["guild_id"],
	})
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_COMBAT_ACHIEVEMENT_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Complete a combat achievement
// @Description	Award points to all users and mark new completers
// @Tags			CombatAchievement
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string									true	"Guild ID"
// @Param			ca_name		path		string									true	"Combat Achievement Name"
// @Param			body		body		models.CompleteCombatAchievementBody		true	"Users"
// @Success		200			{object}	[]database.UpdatePointsByEventRow
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/combat-achievements/{ca_name}/complete [POST]
func (s *Server) CompleteCombatAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	var body models.CompleteCombatAchievementBody
	if err := utils.ParseAndValidateRequestBody(w, r, &body); err != nil {
		return
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	q := s.queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	// 1. Validate CA exists and get point_source
	ca, err := q.GetCombatAchievement(r.Context(), database.GetCombatAchievementParams{
		Name:    p["ca_name"],
		GuildID: p["guild_id"],
	})
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	// 2. Award points to all users
	points, err := q.UpdatePointsByEvent(r.Context(), database.UpdatePointsByEventParams{
		Event:   ca.PointSource,
		GuildID: p["guild_id"],
		UserIds: body.UserIds,
	})
	ei = database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	// 3. Mark all users as having completed the CA (duplicates ignored)
	err = q.CompleteCombatAchievement(r.Context(), database.CompleteCombatAchievementParams{
		UserIds:               body.UserIds,
		GuildID:               p["guild_id"],
		CombatAchievementName: p["ca_name"],
	})
	ei = database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if err = tx.Commit(r.Context()); err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	jw.WriteResponse(points)
}

// @Summary		Grant a combat achievement to a user
// @Description	Mark a user as having completed a combat achievement (no points awarded)
// @Tags			CombatAchievement
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			user_id		path		string	true	"User ID"
// @Param			ca_name		path		string	true	"Combat Achievement Name"
// @Success		201			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/combat-achievements/{ca_name} [POST]
func (s *Server) GiveUserCombatAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)
	p := mux.Vars(r)

	ei := database.WrapExec(s.queries.GiveUserCombatAchievement, r.Context(), database.GiveUserCombatAchievementParams{
		UserID:                p["user_id"],
		GuildID:               p["guild_id"],
		CombatAchievementName: p["ca_name"],
	})
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Remove a combat achievement from a user
// @Description	Remove a user's combat achievement completion record
// @Tags			CombatAchievement
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			user_id		path		string	true	"User ID"
// @Param			ca_name		path		string	true	"Combat Achievement Name"
// @Success		204			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/combat-achievements/{ca_name} [DELETE]
func (s *Server) RemoveUserCombatAchievement(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)
	p := mux.Vars(r)

	rows, err := s.queries.RemoveUserCombatAchievement(r.Context(), database.RemoveUserCombatAchievementParams{
		UserID:                p["user_id"],
		GuildID:               p["guild_id"],
		CombatAchievementName: p["ca_name"],
	})
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_COMBAT_ACHIEVEMENT_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}
