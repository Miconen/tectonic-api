package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary		Get all guild times
// @Description	Get all guild times in a detailed way
// @Tags			Guild
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		200			{object}	models.GuildTimes
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/times [GET]
func (s *Server) GetGuildTimes(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	row, err := s.queries.GetDetailedGuild(r.Context(), guildId)
	ei := database.ClassifyError(err)
	if err != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	guild := database.NewDetailedGuildFromRow(row)
	jw.WriteResponse(guild)
}

// @Summary		Add a new best time to guild
// @Description	Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags			Time
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string				true	"Guild ID"
// @Param			time		body		models.InputTime	true	"Time"
// @Success		200			{object}	models.Empty
// @Success		201			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		409			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/times [POST]
func (s *Server) CreateTime(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)
	var body models.InputTime

	if err := utils.ParseAndValidateRequestBody(w, r, &body); err != nil {
		return
	}

	res := models.TimeResponse{
		BossName: body.BossName,
		Time:     body.Time,
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	q := s.queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	pb_params := database.CheckPbParams{
		Boss:    body.BossName,
		GuildID: p["guild_id"],
	}

	pb, err := q.CheckPb(r.Context(), pb_params)
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	// Old pb exists
	if pb.Time.Valid {
		old_time := int(pb.Time.Int32)
		// Check if our time is faster, if not don't continue
		if old_time <= body.Time {
			jw.WriteResponse(res)
			return
		}
	}

	time_params := database.CreateTimeParams{
		Time:     int32(body.Time),
		BossName: body.BossName,
		Date:     pgtype.Timestamp{Time: time.Now(), Valid: true},
		GuildID:  p["guild_id"],
	}

	run_id, err := q.CreateTime(r.Context(), time_params)
	ei = database.ClassifyError(err)
	if err != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	changed_pb_params := database.UpdatePbParams{
		RunID: pgtype.Int4{
			Int32: run_id,
			Valid: true,
		},
		GuildID: p["guild_id"],
		Boss:    body.BossName,
	}

	_, err = q.UpdatePb(r.Context(), changed_pb_params)
	ei = database.ClassifyError(err)
	if err != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	team_params := database.CreateTeamParams{
		RunID:   run_id,
		UserIds: body.UserIds,
		GuildID: p["guild_id"],
	}

	err = q.CreateTeam(r.Context(), team_params)
	ei = database.ClassifyError(err)
	if err != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	err = tx.Commit(r.Context())
	ei = database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	jw.SetStatus(http.StatusCreated)
	res.RunID = int(run_id)
	res.OldTime = int(pb.Time.Int32)
	jw.WriteResponse(res)
}

// @Summary		Remove time from guilds best times
// @Description	Delete a time in our backend by unique guild Snowflake (ID)
// @Tags			Time
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			time_id		path		string	true	"Time ID"
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/times/{time_id} [DELETE]
func (s *Server) RemoveTime(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteTimeParams{
		GuildID: p["guild_id"],
	}

	id, err := strconv.Atoi(p["time_id"])
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	params.RunID = int32(id)

	deleted, err := s.queries.DeleteTime(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if deleted == 0 {
		jw.WriteError(models.ERROR_TIME_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}
