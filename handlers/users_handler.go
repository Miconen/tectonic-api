package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

//	@Summary		Get one or more users by ID(s)
//	@Description	Get user details by unique user Snowflake (ID)
//	@Tags			Users
//	@Produce		json
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			user_ids	path		string	true	"User ID(s)"
//	@Success		200			{object}	database.DetailedUser[]
//	@Failure		400			{object}	models.Empty
//	@Failure		401			{object}	models.Empty
//	@Failure		404			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/{user_ids} [GET]
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetDetailedUsersParams{
		GuildID: p["guild_id"],
		UserIds: strings.Split(p["user_ids"], ","),
	}

	rows, err := queries.GetDetailedUsers(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if len(rows) == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times, Events: row.Events}
		users = append(users, user)
	}

	jw.WriteResponse(users)
}

//	@Summary		Get one or more users by RSN(s)
//	@Description	Get user details by unique user Snowflake (ID)
//	@Tags			Users
//	@Produce		json
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			rsns		path		string	true	"User RSN(s)"
//	@Success		200			{object}	database.User[]
//	@Failure		400			{object}	models.Empty
//	@Failure		401			{object}	models.Empty
//	@Failure		404			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/rsn/{rsns} [GET]
func GetUsersByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetDetailedUsersByRSNParams{
		GuildID: p["guild_id"],
		Rsns:    strings.Split(p["rsns"], ","),
	}

	rows, err := queries.GetDetailedUsersByRSN(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if len(rows) == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times, Events: row.Events}
		users = append(users, user)
	}

	jw.WriteResponse(users)
}

//	@Summary		Get one or more users by WomID(s)
//	@Description	Get user details by unique user Snowflake (ID)
//	@Tags			Users
//	@Produce		json
//	@Param			guild_id	path		string		true	"Guild ID"
//	@Param			wom_ids		path		[]string	true	"User WomID(s)"
//	@Success		200			{object}	database.User[]
//	@Failure		400			{object}	models.Empty
//	@Failure		401			{object}	models.Empty
//	@Failure		404			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/wom/{wom_ids} [GET]
func GetUsersByWom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetDetailedUsersByWomIDParams{
		GuildID: p["guild_id"],
		WomIds:  strings.Split(p["wom_ids"], ","),
	}

	rows, err := queries.GetDetailedUsersByWomID(r.Context(), params)
	ei :=  database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if len(rows) == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times, Events: row.Events}
		users = append(users, user)
	}

	jw.WriteResponse(users)
}

//	@Summary		Create / Initialize a new user
//	@Description	Initialize a user in our backend by unique user Snowflake (ID)
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			user_id		path		string	true	"User ID"
//	@Param			rsn			body		string	true	"RSN"
//	@Success		201			{object}	models.Empty
//	@Failure		400			{object}	models.ErrorResponse
//	@Failure		401			{object}	models.Empty
//	@Failure		409			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/{user_id} [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)

	v := mux.Vars(r)
	params := database.CreateUserParams{
		GuildID: v["guild_id"],
		UserID:  v["user_id"],
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_BODY)
		return
	}

	wom, err := utils.GetWom(params.Rsn)
	if err != nil {
		jw.WriteError(models.ERROR_RSN_NOT_FOUND)
		return
	}

	params.WomID = strconv.Itoa(wom.Id)
	params.Rsn = wom.DisplayName

	user, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		ei := database.ClassifyError(err)
		if ei != nil {
			handleDatabaseErrorCustom(*ei, jw, func(dh *dbHandler, jw *utils.JsonWriter) {
				switch dh.Err.ConstraintName {
				case "users_ibfk_1":
					jw.WriteError(models.ERROR_GUILD_EXISTS)
				case "users_pkey":
					jw.WriteError(models.ERROR_USER_EXISTS)
				}
			})
			return
		}
	}

	jw.WriteResponse(user)
}

//	@Summary		Delete a user from guild by User ID
//	@Description	Delete a user in our backend by unique user and guild Snowflake (ID)
//	@Tags			User
//	@Produce		json
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			user_id		path		string	true	"User ID"
//	@Success		204			{object}	models.Empty
//	@Failure		400			{object}	models.Empty
//	@Failure		401			{object}	models.Empty
//	@Failure		404			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/{user_id} [DELETE]
func RemoveUserById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteUserByIdParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
	}

	rows, err := queries.DeleteUserById(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
	}

	jw.WriteResponse(http.NoBody)
}

//	@Summary		Delete a user from guild by RSN
//	@Description	Delete a user in our backend by Runescape name
//	@Tags			User
//	@Produce		json
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			rsn			path		string	true	"RSN"
//	@Success		204			{object}	models.Empty
//	@Failure		400			{object}	models.Empty
//	@Failure		401			{object}	models.Empty
//	@Failure		404			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/rsn/{rsn} [DELETE]
func RemoveUserByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteUserByRsnParams{
		GuildID: p["guild_id"],
		Rsn:     p["rsn"],
	}

	rows, err := queries.DeleteUserByRsn(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
	}

	jw.WriteResponse(http.NoBody)
}

//	@Summary		Delete a user from guild by Wom ID
//	@Description	Delete a user in our backend by unique user and guild Snowflake (ID)
//	@Tags			User
//	@Produce		json
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			wom_id		path		string	true	"Wom ID"
//	@Success		204			{object}	models.Empty
//	@Failure		400			{object}	models.Empty
//	@Failure		401			{object}	models.Empty
//	@Failure		404			{object}	models.Empty
//	@Failure		429			{object}	models.Empty
//	@Failure		500			{object}	models.Empty
//	@Router			/api/v1/guilds/{guild_id}/users/wom/{wom_id} [DELETE]
func RemoveUserByWom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	params := database.DeleteUserByWomParams{
		GuildID: p["guild_id"],
		WomID:   p["wom_id"],
	}

	rows, err := queries.DeleteUserByWom(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
	}

	jw.WriteResponse(http.NoBody)
}
