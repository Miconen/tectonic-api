package handlers

import (
	"encoding/json"
	"net/http"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary Get a guild by ID
// @Description Get guild details by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} models.Guild
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id} [GET]
func GetGuild(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		log.Error("no guild id found")
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(http.NoBody)
		return
	}

	guild, err := queries.GetGuild(r.Context(), guildId)
	if err != nil {
		log.Error("Error selecting guild", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	// Write JSON response
	jw.WriteResponse(guild)
}

// @Summary Create / Initialize a guild
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Accept json
// @Produce json
// @Param guild body models.InputGuild true "Guild"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds [POST]
func CreateGuild(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)

	p := models.InputGuild{
		Multiplier: 1,
	}

	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	_, err = queries.CreateGuild(r.Context(), p.GuildId)
	if err != nil {
		log.Error("Error creating guild", "error", err)
		jw.SetStatus(http.StatusConflict)
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary Delete a guild
// @Description Delete a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id} [DELETE]
func DeleteGuild(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		log.Error("no guild id found")
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(http.NoBody)
		return
	}

	_, err := queries.DeleteGuild(r.Context(), guildId)
	if err != nil {
		log.Error("Error deleting guild", "error", err)
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(http.NoBody)
}

type CategoryMessage struct {
	MessageID string `json:"message_id"`
	Category  string `json:"category"`
}

type GuildParams struct {
	Multiplier       pgtype.Numeric    `json:"multiplier"`
	PbChannelID      string            `json:"pb_channel_id"`
	CategoryMessages []CategoryMessage `json:"category_messages"`
}

// @Summary Updates a guild
// @Description Update multiplier and/or time channel for a guild
// @Tags Guild
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param guild body models.UpdateGuild true "Guild"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id} [PUT]
func UpdateGuild(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusNoContent)

	p := mux.Vars(r)

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		log.Error("Error creating transaction", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
	}

	q := queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	var params GuildParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Error("Failed to parse request body")
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	categories := make([]string, len(params.CategoryMessages))
	messageIds := make([]string, len(params.CategoryMessages))
	for i, v := range params.CategoryMessages {
		categories[i] = v.Category
		messageIds[i] = v.MessageID
	}

	if len(params.CategoryMessages) > 0 {
		message_params := database.UpdateCategoryMessageIdsParams{
			GuildID:    p["guild_id"],
			Categories: categories,
			MessageIds: messageIds,
		}

		_, errUpdate := q.UpdateCategoryMessageIds(r.Context(), message_params)
		if errUpdate != nil {
			log.Error("Error updating channel", "error", errUpdate)
			jw.SetStatus(http.StatusNotFound)
		}
	}

	guild_params := database.UpdateGuildParams{
		Multiplier:  params.Multiplier,
		PbChannelID: params.PbChannelID,
		GuildID:     p["guild_id"],
	}
	
	_, err = q.UpdateGuild(r.Context(), guild_params)
	if err != nil {
		log.Error("Error updating channel", "error", err)
		jw.SetStatus(http.StatusNotFound)
	}

	tx.Commit(r.Context())

	jw.WriteResponse(http.NoBody)
}
