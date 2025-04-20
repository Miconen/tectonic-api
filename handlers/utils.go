package handlers

import (
	"context"
	"net/http"
	"regexp"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbHandler struct {
	database.ErrorInfo
}

// You can get this by running the following SQL query:
// SELECT constraint_name, table_name FROM information_schema.table_constraints WHERE table_schema = 'public'
var constraintsMap map[string]database.ConstraintDetail

func InitDatabaseHandler(ctx context.Context, conn *pgxpool.Conn) error {
	var err error
	constraintsMap, err = database.GetConstraintsTable(ctx, conn)
	return err
}

func getConstraintError(ei database.ErrorInfo) models.APIV1Error {
	switch ei.Code {
	case "23505":
		c := constraintsMap[ei.Err.ConstraintName]

		switch c.Type {
		case database.PrimaryKey:
			switch c.Table {
			case "achievement":
				return models.ERROR_ACHIEVEMENT_EXISTS
			case "bosses":
				return models.ERROR_BOSS_EXISTS
			case "categories":
				return models.ERROR_CATEGORY_EXISTS
			case "event":
				return models.ERROR_EVENT_EXISTS
			case "event_participant":
				return models.ERROR_PARTICIPATION_EXISTS
			case "guild_bosses":
				return models.ERROR_GUILD_BOSS_EXISTS
			case "guild_categories":
				return models.ERROR_GUILD_CATEGORY_EXISTS
			case "guilds":
				return models.ERROR_GUILD_EXISTS
			case "point_sources":
				return models.ERROR_POINT_SOURCE_EXISTS
			case "rsn":
				return models.ERROR_RSN_EXISTS
			case "teams":
				return models.ERROR_TEAM_EXISTS
			case "times":
				return models.ERROR_TIME_EXISTS
			case "user_achievement":
				return models.ERROR_USER_ACHIEVEMENT_EXISTS
			case "users":
				return models.ERROR_USER_EXISTS
			}
		case database.ForeignKey:
			switch c.ForeignTable {
			case "achievement":
				return models.ERROR_ACHIEVEMENT_NOT_FOUND
			case "bosses":
				return models.ERROR_BOSS_NOT_FOUND
			case "categories":
				return models.ERROR_CATEGORY_NOT_FOUND
			case "event":
				return models.ERROR_EVENT_NOT_FOUND
			case "event_participant":
				return models.ERROR_PARTICIPATION_NOT_FOUND
			case "guild_bosses":
				return models.ERROR_GUILD_BOSS_NOT_FOUND
			case "guild_categories":
				return models.ERROR_GUILD_CATEGORY_NOT_FOUND
			case "guilds":
				return models.ERROR_GUILD_NOT_FOUND
			case "point_sources":
				return models.ERROR_POINT_SOURCE_NOT_FOUND
			case "rsn":
				return models.ERROR_RSN_NOT_FOUND
			case "teams":
				return models.ERROR_TEAM_NOT_FOUND
			case "times":
				return models.ERROR_TIME_NOT_FOUND
			case "user_achievement":
				return models.ERROR_USER_ACHIEVEMENT_NOT_FOUND
			case "users":
				return models.ERROR_USER_NOT_FOUND
			}
		}
	}

	log.Error("Database untreated error", "error_info", ei)
	return models.ERROR_TODO(400, time.Now().String())
}

// Handlers should delegate database errors to this function always.
// This handler will alwas write the response, so the caller should
// always short circuit the handler.
func handleDatabaseError(ei database.ErrorInfo, jw *utils.JsonWriter) {
	dh := &dbHandler{ErrorInfo: ei}

	if dh.Severity == database.SeverityFatal || dh.Severity == database.SeverityPanic {
		log.Error("DATABASE FAILURE", "err", dh.Error())
		jw.WriteError(models.ERROR_API_DEAD)
	} else if !dh.Recoverable {
		log.Error("DATABASE NON-RECOVERABLE ERROR", "err", dh.Error())
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
	} else {
		log.Warn("QUERY ERROR", "err", dh.Error())
		jw.WriteError(getConstraintError(ei))
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
		log.Warn("QUERY ERROR", "err", dh.Error())
		dhc(dh, jw)
	}
}

// Middleware that validate URL parameters for the handler. POST requests are
// NOT validated.
func ValidateParameters(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jw := utils.NewJsonWriter(w, r, 0)
		p := mux.Vars(r)

		// URL parameter needs to be excluded from validation because it will
		// get created later on the request.
		if r.Method == "POST" {
			pat := regexp.MustCompile(`{\w*}`)
			vars := pat.FindAllString(r.URL.String(), -1)

			if len(vars) != 0 {
				switch vars[len(vars)-1] {
				case "guild_id", "user_id", "event_id":
					delete(p, vars[len(vars)-1])
				}
			}
		}

		ctx := r.Context()
		conn, err := database.AcquireConnection(ctx)
		if err != nil {
			jw.WriteError(models.ERROR_API_UNAVAILABLE)
			return
		}
		defer conn.Release()

		for k, v := range p {
			var ok bool
			switch k {
			case "guild_id":
				ok = guildExists(ctx, conn, jw, v)
			case "user_id":
				ok = userExists(ctx, conn, jw, v)
			case "event_id":
				ok = eventExists(ctx, conn, jw, v)
			case "achievement_name":
				ok = achievementExists(ctx, conn, jw, v)
			}

			if !ok {
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

// Checks if guild exists on the database.
func guildExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, guild_id string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT guild_id FROM guilds WHERE guild_id = $1)", guild_id, models.ERROR_GUILD_NOT_FOUND)
}

// Checks if user exists on the database.
func userExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, user_id string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT user_id FROM users WHERE user_id = $1)", user_id, models.ERROR_USER_NOT_FOUND)
}

// Checks if event exists on the database.
func eventExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, event_id string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT wom_id FROM event WHERE wom_id = $1)", event_id, models.ERROR_EVENT_NOT_FOUND)
}

// Checks if achievent exists on the database.
func achievementExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, name string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT name FROM achievement WHERE name = $1)", name, models.ERROR_ACHIEVEMENT_NOT_FOUND)
}

func queryExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, sql string, param string, api_err models.APIV1ErrorCode) bool {
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
