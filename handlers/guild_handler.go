package handlers

import "net/http"

// @Summary Get a guild by ID
// @Description Get guild details by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Success 200 {object} Guild
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/guild [GET]
func GetGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		guild := Guild{
			GuildId:     p["guild_id"],
			Multiplier:  1,
			PbChannelId: "123",
		}

		return guild, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Create / Initialize a guild
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} Guild
// @Success 201 {object} Guild
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/guild [POST]
func CreateGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		guild := Guild{
			GuildId:     p["guild_id"],
			Multiplier:  1,
			PbChannelId: "123",
		}

		return guild, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Update a guilds information
// @Description Update a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} Guild
// @Success 201 {object} Guild
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/guild [PUT]
func UpdateGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		guild := Guild{
			GuildId:     p["guild_id"],
			Multiplier:  1,
			PbChannelId: "123",
		}

		return guild, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Delete a guild
// @Description Delete a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} Guild
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/guild [DELETE]
func RemoveGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		guild := Guild{
			GuildId:     p["guild_id"],
			Multiplier:  1,
			PbChannelId: "123",
		}

		return guild, nil
	}

	httpHandler(w, r, h, p)
}
