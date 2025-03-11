package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
)

// @Summary Get one or more users by ID(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_ids path string true "User ID(s)"
// @Success 200 {object} database.DetailedUSer[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_ids} [GET]
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetDetailedUsersParams{
		GuildID: p["guild_id"],
		UserIds: strings.Split(p["user_ids"], ","),
	}

	rows, err := queries.GetDetailedUsers(r.Context(), params)
	if err != nil {
		log.Error("Error fetching user", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times}
		users = append(users, user)
	}

	jw.WriteResponse(users)
}

// @Summary Get one or more users by RSN(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param rsns path string true "User RSN(s)"
// @Success 200 {object} database.User[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/rsn/{rsns} [GET]
func GetUsersByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetDetailedUsersByRSNParams{
		GuildID: p["guild_id"],
		Rsns:    strings.Split(p["rsns"], ","),
	}

	rows, err := queries.GetDetailedUsersByRSN(r.Context(), params)
	if err != nil {
		log.Error("Error fetching user", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times}
		users = append(users, user)
	}

	jw.WriteResponse(users)
}

// @Summary Get one or more users by WomID(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param wom_id path string true "User WomID(s)"
// @Success 200 {object} database.User[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/wom/{wom_ids} [GET]
func GetUsersByWom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetDetailedUsersByWomIDParams{
		GuildID: p["guild_id"],
		WomIds:  strings.Split(p["wom_ids"], ","),
	}

	rows, err := queries.GetDetailedUsersByWomID(r.Context(), params)
	if err != nil {
		log.Error("Error fetching user", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times}
		users = append(users, user)
	}

	jw.WriteResponse(users)
}

// @Summary Create / Initialize a new user
// @Description Initialize a user in our backend by unique user Snowflake (ID)
// @Tags User
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn body string true "RSN"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id} [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)

	v := mux.Vars(r)
	params := database.CreateUserParams{
		GuildID: v["guild_id"],
		UserID:  v["user_id"],
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Error("Error decoding request body: ", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(err)
		return
	}

	wom, err := utils.GetWom(params.Rsn)
	if err != nil {
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	params.WomID = strconv.Itoa(wom.Id)
	params.Rsn = wom.DisplayName

	user, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		log.Error("Error inserting user", "error", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_ibfk_1":
				jw.SetStatus(http.StatusBadRequest)
			case "users_pkey":
				jw.SetStatus(http.StatusConflict)
			default:
				jw.SetStatus(http.StatusInternalServerError)
				jw.WriteResponse(http.NoBody)
				return
			}
		}
	}

	jw.WriteResponse(user)
}

// @Summary Delete a user from guild by User ID
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id} [DELETE]
func RemoveUserById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteUserByIdParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
	}

	rows, err := queries.DeleteUserById(r.Context(), params)
	if err != nil {
		log.Error("Error deleting user", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
	}

	if rows == 0 {
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary Delete a user from guild by RSN
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param rsn path string true "RSN"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/rsn/{rsn} [DELETE]
func RemoveUserByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteUserByRsnParams{
		GuildID: p["guild_id"],
		Rsn:     p["rsn"],
	}

	rows, err := queries.DeleteUserByRsn(r.Context(), params)
	if err != nil {
		log.Error("Error deleting user", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
	}

	if rows == 0 {
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary Delete a user from guild by Wom ID
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param wom_id path string true "Wom ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/wom/{wom_id} [DELETE]
func RemoveUserByWom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteUserByWomParams{
		GuildID: p["guild_id"],
		WomID:   p["wom_id"],
	}

	rows, err := queries.DeleteUserByWom(r.Context(), params)
	if err != nil {
		log.Error("Error deleting user", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
	}

	if rows == 0 {
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(http.NoBody)
}
