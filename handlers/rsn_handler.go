package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
)

// @Summary Link an RSN to a user
// @Description Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn body models.InputRSN true "RSN"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id}/rsns/{rsn} [POST]
func CreateRSN(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)
	params := database.CreateRsnParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
		Rsn:     p["rsn"],
	}

	wom, err := utils.GetWom(params.Rsn)
	if err != nil {
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	params.WomID = strconv.Itoa(wom.Id)
	params.Rsn = wom.DisplayName

	err = queries.CreateRsn(r.Context(), params)
	if err != nil {
		log.Error("Error creating RSN", "error", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				// Foreign key violation (User not found)
				jw.SetStatus(http.StatusNotFound)
			} else if pgErr.Code == "23505" {
				// Unique violation (Duplicate)
				jw.SetStatus(http.StatusConflict)
			} else {
				jw.SetStatus(http.StatusInternalServerError)
			}
		}
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary Remove RSN from guild and user
// @Description Delete a RSN in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn path string true "RSN"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id}/rsns/{rsn} [DELETE]
func RemoveRSN(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)
	params := database.DeleteRsnParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
		Rsn:     p["rsn"],
	}

	rows, err := queries.DeleteRsn(r.Context(), params)
	if err != nil {
		log.Error("Error deleting RSN", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
	}

	if rows == 0 {
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(http.NoBody)
}
