package handlers

import (
	"context"
	"net/http"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
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

// Generic function validation type
type existsFunc func(ctx context.Context, jw *utils.JsonWriter, conn *pgxpool.Conn) bool

// Middleware that validate URL parameters for the handler.
func ValidateParameters(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jw := utils.NewJsonWriter(w, r, 0)
		funcs := []existsFunc{}

		p := mux.Vars(r)
		for k, v := range p {
			switch k {
			case "guild_id":
				funcs = append(funcs, guildExists(v))
			case "user_id":
				funcs = append(funcs, userExists(v))
			}
		}

		if !exists(r.Context(), jw, funcs...) {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Validation function aggregate that serves to validate request parameters.
// Returns true if all of them are valid, otherwise false.
func exists(ctx context.Context, jw *utils.JsonWriter, funcs ...existsFunc) bool {
	conn, err := database.AcquireConnection(ctx)
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return false
	}
	defer conn.Release()

	for _, f := range funcs {
		if !f(ctx, jw, conn) {
			return false
		}
	}

	return true
}

// Checks if guild exists on the database.
func guildExists(guild_id string) existsFunc {
	return queryExists("SELECT EXISTS (SELECT guild_id FROM guilds WHERE guild_id = $1)", guild_id, models.ERROR_GUILD_NOT_FOUND)
}

// Checks if user exists on the database.
func userExists(user_id string) existsFunc {
	return queryExists("SELECT EXISTS (SELECT user_id FROM users WHERE user_id = $1)", user_id, models.ERROR_USER_NOT_FOUND)
}

func queryExists(sql string, param string, api_err models.APIV1Error) existsFunc {
	return func(ctx context.Context, jw *utils.JsonWriter, conn *pgxpool.Conn) bool {
		exists := false
		err := conn.QueryRow(ctx, sql, param).Scan(&exists)
		if err != nil {
			jw.WriteError(models.ERROR_API_DEAD)
			return false
		}

		if !exists {
			jw.WriteError(api_err)
			return false
		}

		return true
	}
}
