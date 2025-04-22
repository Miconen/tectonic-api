package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary		Link an RSN to a user
// @Description	Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags			RSN
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string			true	"Guild ID"
// @Param			user_id		path		string			true	"User ID"
// @Param			user_id		body		string			true	"User ID"
// @Param			rsn			path		models.InputRSN	true	"RSN"
// @Success		204			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		409			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/rsns [POST]
func CreateRSN(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)
	params := database.CreateRsnParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
	}

	err := utils.ParseRequestBody(w, r, &params)
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_BODY)
		return
	}

	wom, err := utils.GetWom(params.Rsn)
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	params.WomID = strconv.Itoa(wom.Id)
	params.Rsn = wom.DisplayName

	err = queries.CreateRsn(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseErrorCustom(*ei, jw, func(dh *dbHandler, jw *utils.JsonWriter) {
			switch dh.Code {
			case "23503":
				jw.WriteResponse(models.ERROR_USER_NOT_FOUND)
			case "23505":
				jw.WriteResponse(models.ERROR_GUILD_EXISTS)
			}
		})
		return
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Remove RSN from guild and user
// @Description	Delete a RSN in our backend by unique guild and user Snowflake (ID)
// @Tags			RSN
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Param			user_id		path		string	true	"User ID"
// @Param			rsn			path		string	true	"RSN"
// @Success		201			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		404			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_id}/rsns/{rsn} [DELETE]
func RemoveRSN(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)
	params := database.DeleteRsnParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
		Rsn:     p["rsn"],
	}

	rows, err := queries.DeleteRsn(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw)
		return
	}

	if rows == 0 {
		jw.WriteError(models.ERROR_RSN_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}
