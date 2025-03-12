package handlers

import (
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"
)

type dbHandler struct {
	database.ErrorInfo
}

// You can get this by running the following SQL query:
// SELECT constraint_name, table_name FROM information_schema.table_constraints WHERE table_schema = 'public'
var constraintsMap = map[string]string{
	"categories_name":                    "Categories",
	"boss_name":                          "Boss",
	"guild_bosses_bosses_guild_id":       "Guild",
	"guild_categories_guild_id_category": "Guild",
	"guilds_pkey":                        "Guilds",
	"rsn_pkey":                           "Rsn",
	"teams_run_id_user_id_guild_id":      "Teams",
	"times_pkey":                         "Times",
	"users_pkey":                         "Users",
	"point_sources_pkey":                 "Point",
	"boss_category_fkey":                 "Boss",
	"guild_bosses_bosses_fkey":           "Guild",
	"guild_bosses_guild_id_fkey":         "Guild",
	"guild_bosses_pb_id_fkey":            "Guild",
	"guild_categories_category_fkey":     "Guild",
	"guild_categories_guild_id_fkey":     "Guild",
	"rsn_ibfk_1":                         "Rsn",
	"teams_run_id_fkey":                  "Teams",
	"teams_user_id_guild_id_fkey":        "Teams",
	"times_bosses_name_fkey":             "Times",
	"users_ibfk_1":                       "Users",
	"point_sources_ibfk_1":               "Point",
	"times_guild_id_fkey":                "Times",
}

// Handlers should delegate database errors to this handler always, and check
// the error if the handler wrote the response. This handler will only write
// the response if the application can't recover the error, passing them to
// the client instead.
func handleDatabaseError(ei database.ErrorInfo, jw *utils.JsonWriter, code models.APIV1Error) {
	dh := &dbHandler{ErrorInfo: ei}

	if dh.Severity == database.SeverityFatal || dh.Severity == database.SeverityPanic {
		log.Error("DATABASE FAILURE", "err", dh.Error())
		jw.WriteError(models.ERROR_API_DEAD)
	} else if !dh.Recoverable {
		log.Error("DATABASE NON-RECOVERABLE ERROR", "err", dh.Error())
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
	} else {
		log.Warn("DATABASE ERROR", "err", dh.Error())
		jw.WriteError(code)
	}
}

type dbHandlerFunc func(dh *dbHandler, jw *utils.JsonWriter)

func handleDatabaseErrorCustom(ei database.ErrorInfo, jw *utils.JsonWriter, dhc dbHandlerFunc) {
	dh := &dbHandler{ErrorInfo: ei}

	if dh.Severity == database.SeverityFatal || dh.Severity == database.SeverityPanic {
		log.Error("DATABASE FAILURE", "err", dh.Error())
		jw.WriteError(models.ERROR_API_DEAD)
	} else if !dh.Recoverable {
		log.Error("DATABASE NON-RECOVERABLE ERROR", "err", dh.Error())
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
	} else {
		log.Warn("DATABASE ERROR", "err", dh.Error())
		dhc(dh, jw)
	}
}
