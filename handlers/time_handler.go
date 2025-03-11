package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary Get all guild times
// @Description Get all guild times in a detailed way
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} models.GuildTimes
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/times [GET]
func GetGuildTimes(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		log.Error("no guild id found")
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(http.NoBody)
		return
	}

	row, err := queries.GetDetailedGuild(r.Context(), guildId)
	if err != nil {
		log.Error("Error selecting guild", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	guild := database.NewDetailedGuildFromRow(row)
	jw.WriteResponse(guild)
}

// @Summary Add a new best time to guild
// @Description Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param time body models.InputTime true "Time"
// @Success 200 {object} models.Empty
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guilds/{guild_id}/times [POST]
func CreateTime(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := models.InputTime{}
	err := utils.ParseRequestBody(w, r, &params)
	if err != nil {
		log.Error("Error parsing request body", "error", err)
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	if len(params.UserIds) == 0 {
		log.Error("Empty User IDs array not permitted")
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	res := models.TimeResponse{
		BossName: params.BossName,
		Time:     params.Time,
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		log.Error("Error creating transaction", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	q := queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	pb_params := database.CheckPbParams{
		Boss:    params.BossName,
		GuildID: p["guild_id"],
	}

	pb, err := q.CheckPb(r.Context(), pb_params)
	if err != nil {
		log.Error("Error checking pb", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	// Old pb exists
	if pb.Time.Valid {
		old_time := int(pb.Time.Int32)
		// Check if our time is faster, if not don't continue
		if old_time <= params.Time {
			jw.WriteResponse(res)
			return
		}
	}

	time_params := database.CreateTimeParams{
		Time:     int32(params.Time),
		BossName: params.BossName,
		Date:     pgtype.Timestamp{Time: time.Now(), Valid: true},
		GuildID:  p["guild_id"],
	}

	run_id, err := q.CreateTime(r.Context(), time_params)
	if err != nil {
		log.Error("Error inserting time", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	changed_pb_params := database.UpdatePbParams{
		RunID: pgtype.Int4{
			Int32: run_id,
			Valid: true,
		},
		GuildID: p["guild_id"],
		Boss:    params.BossName,
	}

	_, err = q.UpdatePb(r.Context(), changed_pb_params)
	if err != nil {
		log.Error("Error updating pb", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	team_params := database.CreateTeamParams{
		RunID:   run_id,
		UserIds: params.UserIds,
		GuildID: p["guild_id"],
	}

	err = q.CreateTeam(r.Context(), team_params)
	if err != nil {
		log.Error("Error creating team", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	tx.Commit(r.Context())

	jw.SetStatus(http.StatusCreated)
	res.RunID = int(run_id)
	res.OldTime = int(pb.Time.Int32)
	jw.WriteResponse(res)
}

// @Summary Remove time from guilds best times
// @Description Delete a time in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param time_id path string true "Time ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guilds/{guild_id}/times/{time_id} [DELETE]
func RemoveTime(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteTimeParams{
		GuildID: p["guild_id"],
	}

	id, err := strconv.Atoi(p["time_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params.RunID = int32(id)

	deleted, err := queries.DeleteTime(r.Context(), params)
	if err != nil {
		log.Error("Error deleting time", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	if deleted == 0 {
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	jw.WriteResponse(http.NoBody)
}
