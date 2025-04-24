package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

func getDetailedUsers(ctx context.Context, user_ids []string, guild_id string) ([]models.DetailedUser, *database.ErrorInfo) {
	detailed_users := make([]models.DetailedUser, len(user_ids))

	for i, user_id := range user_ids {
		user_rows, err := database.WrapQuery(queries.GetUsersById, ctx, database.GetUsersByIdParams{
			GuildID: guild_id,
			UserIds: []string{user_id},
		})
		if err != nil {
			return nil, err
		}

		if len(user_rows) != 1 {
			return nil, err
		}
		user := user_rows[0]

		times_rows, err := database.WrapQuery(queries.GetUserTimes, ctx, database.GetUserTimesParams{
			UserID:  user_id,
			GuildID: guild_id,
		})
		if err != nil {
			return nil, err
		}

		achievements_rows, err := database.WrapQuery(queries.GetUserAchievements, ctx, user_id)
		if err != nil {
			return nil, err
		}

		events_rows, err := database.WrapQuery(queries.GetUserEvents, ctx, database.GetUserEventsParams{
			UserID:  user_id,
			GuildID: guild_id,
		})
		if err != nil {
			return nil, err
		}

		detailed_users[i] = models.DetailedUser{
			UserId:       user.UserID,
			GuildId:      user.GuildID,
			Points:       int(user.Points),
			Times:        models.UserTimesFromRows(times_rows),
			Events:       models.UserEventFromRows(events_rows),
			Achievements: models.UserAchievementsFromRows(achievements_rows),
		}
	}

	return detailed_users, nil
}

// @Summary		Get one or more users by ID(s)
// @Description	Get user details by unique user Snowflake (ID)
// @Tags			Users
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			user_ids	path		string	true	"User ID(s)"
// @Success		200			{object}	database.DetailedUser[]
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_ids} [GET]
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	detailed_users, err := getDetailedUsers(r.Context(), strings.Split(p["user_ids"], ","), p["guild_id"])
	if err != nil {
		handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(detailed_users)
}

// @Summary		Get one or more users by RSN(s)
// @Description	Get user details by unique user Snowflake (ID)
// @Tags			Users
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			rsns		path		string	true	"User RSN(s)"
// @Success		200			{object}	database.User[]
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/rsn/{rsns} [GET]
func GetUsersByRsn(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	users, err := database.WrapQuery(queries.GetUsersByRsn, r.Context(), database.GetUsersByRsnParams{
		GuildID: p["guild_id"],
		Rsns:    strings.Split(p["rsns"], ","),
	})

	if err != nil {
		handleDatabaseError(*err, jw)
		return
	}

	getDetailedUsers(
		r.Context(),
		utils.MapField(users, func(r database.User) string { return r.UserID }),
		p["guild_id"],
	)
}

// @Summary		Get one or more users by WomID(s)
// @Description	Get user details by unique user Snowflake (ID)
// @Tags			Users
// @Produce		json
// @Param			guild_id	path		string		true	"Guild ID"
// @Param			wom_ids		path		[]string	true	"User WomID(s)"
// @Success		200			{object}	database.User[]
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/wom/{wom_ids} [GET]
func GetUsersByWom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	users, err := database.WrapQuery(queries.GetUsersByWom, r.Context(), database.GetUsersByWomParams{
		GuildID: p["guild_id"],
		WomIds:  strings.Split(p["wom_ids"], ","),
	})

	if err != nil {
		handleDatabaseError(*err, jw)
		return
	}

	getDetailedUsers(
		r.Context(),
		utils.MapField(users, func(r database.User) string { return r.UserID }),
		p["guild_id"],
	)
}

// @Summary		Get the user's achievemnts
// @Description	Get all user's achievemnts registered in the API.
// @Tags			Users
// @Produce		json
// @Param			guild_id	path		string		true	"Guild ID"
// @Param			wom_ids		path		[]string	true	"User WomID(s)"
// @Success		200			{object}	database.GetUserAchievemntsRow[]
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/achievements [GET]
func GetUserAchievements(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	achievements, err := database.WrapQuery(queries.GetUserAchievements, r.Context(), p["user_id"])
	if err != nil {
		handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(achievements)
}

// @Summary		Get the user's events
// @Description	Get all user's events that participated.
// @Tags			Users
// @Produce		json
// @Param			guild_id	path		string		true	"Guild ID"
// @Param			wom_ids		path		[]string	true	"User WomID(s)"
// @Success		200			{object}	database.GetUserEventsRow[]
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/events [GET]
func GetUserEvents(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	events, err := database.WrapQuery(queries.GetUserEvents, r.Context(), database.GetUserEventsParams{
		UserID:  p["user_id"],
		GuildID: p["guild_id"],
	})
	if err != nil {
		handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(events)
}

// @Summary		Get the user times
// @Description	Get all user times.
// @Tags			Users
// @Produce		json
// @Param			guild_id	path		string		true	"Guild ID"
// @Param			wom_ids		path		[]string	true	"User WomID(s)"
// @Success		200			{object}	models.UserTimes[]
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/times [GET]
func GetUserTimes(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	rows, err := database.WrapQuery(queries.GetUserTimes, r.Context(), database.GetUserTimesParams{
		UserID:  p["user_id"],
		GuildID: p["guild_id"],
	})
	if err != nil {
		handleDatabaseError(*err, jw)
		return
	}

	times := models.UserTimesFromRows(rows)
	jw.WriteResponse(times)
}

// @Summary		Create / Initialize a new user
// @Description	Initialize a user in our backend by unique user Snowflake (ID)
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			user_id		body		string	true	"User ID"
// @Param			rsn			body		string	true	"RSN"
// @Success		201			{object}	models.Empty
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.Empty
// @Failure		409			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)

	v := mux.Vars(r)
	params := database.CreateUserParams{
		GuildID: v["guild_id"],
	}

	body := models.CreateUserBody{}
	err := utils.ParseRequestBody(w, r, &body)
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_BODY)
		return
	}

	wom, err := utils.GetWom(body.RSN)
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
			handleDatabaseError(*ei, jw)
			return
		}
	}

	jw.WriteResponse(user)
}

// @Summary		Delete a user from guild by User ID
// @Description	Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags			User
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			user_id		path		string	true	"User ID"
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id} [DELETE]
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
		handleDatabaseError(*ei, jw)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Delete a user from guild by RSN
// @Description	Delete a user in our backend by Runescape name
// @Tags			User
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			rsn			path		string	true	"RSN"
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/rsn/{rsn} [DELETE]
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
		handleDatabaseError(*ei, jw)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Delete a user from guild by Wom ID
// @Description	Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags			User
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			wom_id		path		string	true	"Wom ID"
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/wom/{wom_id} [DELETE]
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
		handleDatabaseError(*ei, jw)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_USER_NOT_FOUND)
	}

	jw.WriteResponse(http.NoBody)
}
