package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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
	status := http.StatusNoContent

	p := mux.Vars(r)
	params := database.CreateRsnParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
		Rsn:     p["rsn"],
	}

	wid, err := utils.GetWomId(params.Rsn)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	params.WomID = wid

	err = queries.CreateRsn(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating RSN: %v\n", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				// Foreign key violation (User not found)
				status = http.StatusNotFound
			} else if pgErr.Code == "23505" {
				// Unique violation (Duplicate)
				status = http.StatusConflict
			} else {
				status = http.StatusInternalServerError
			}
		}
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
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
	status := http.StatusNoContent

	p := mux.Vars(r)
	params := database.DeleteRsnParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
		Rsn:     p["rsn"],
	}

	rows, err := queries.DeleteRsn(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting RSN: %v\n", err)
		status = http.StatusInternalServerError
	}

	if rows == 0 {
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
