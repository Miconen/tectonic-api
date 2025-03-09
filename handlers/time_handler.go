package handlers

import (
	"fmt"
	"net/http"
	"os"
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
	status := http.StatusOK

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		fmt.Fprintf(os.Stderr, "No guild id found")
		status = http.StatusBadRequest
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	row, err := queries.GetDetailedGuild(r.Context(), guildId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting guild: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	guild := database.NewDetailedGuildFromRow(row)
	utils.JsonWriter(guild).IntoHTTP(status)(w, r)
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
	status := http.StatusOK

	p := mux.Vars(r)

	params := models.InputTime{}
	err := utils.ParseRequestBody(w, r, &params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing request body: %v\n", err)
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	if len(params.UserIds) == 0 {
		fmt.Fprintf(os.Stderr, "Empty User IDs array not permitted.\n")
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	res := models.TimeResponse{
		BossName: params.BossName,
		Time:     params.Time,
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating transaction: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
	}

	q := queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	pb_params := database.CheckPbParams{
		Boss:    params.BossName,
		GuildID: p["guild_id"],
	}

	pb, err := q.CheckPb(r.Context(), pb_params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking pb: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	// Old pb exists
	if pb.Time.Valid {
		old_time := int(pb.Time.Int32)
		// Check if our time is faster, if not don't continue
		if old_time <= params.Time {
			utils.JsonWriter(res).IntoHTTP(status)(w, r)
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
		fmt.Fprintf(os.Stderr, "Error inserting time: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
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
		fmt.Fprintf(os.Stderr, "Error updating pb: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	team_params := database.CreateTeamParams{
		RunID:   run_id,
		UserIds: params.UserIds,
		GuildID: p["guild_id"],
	}

	err = q.CreateTeam(r.Context(), team_params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating team: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	tx.Commit(r.Context())

	status = http.StatusCreated
	res.RunID = int(run_id)
	res.OldTime = int(pb.Time.Int32)
	utils.JsonWriter(res).IntoHTTP(status)(w, r)
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
	status := http.StatusNoContent

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
		fmt.Fprintf(os.Stderr, "Error deleting time: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	if deleted == 0 {
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
