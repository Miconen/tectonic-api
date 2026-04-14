package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"tectonic-api/config"
	"tectonic-api/database"
	"tectonic-api/handlers"
	"tectonic-api/models"
	"tectonic-api/routes"
	"tectonic-api/utils"

	"github.com/go-chi/chi/v5"
)

type TestVariables struct {
	GuildId         string
	UserId          string
	Rsn             string
	RsnExtra        string
	ChannelId       string
	Multiplier      int
	WomId           string
	EventClassicId  int
	EventTeamId     int
	AchievementName string
}

func (tv TestVariables) RsnEscaped() string {
	return url.PathEscape(tv.Rsn)
}

func (tv TestVariables) RsnExtraEscaped() string {
	return url.PathEscape(tv.RsnExtra)
}

type TestTable struct {
	Name       string
	Method     string
	Path       string
	Body       any
	StatusCode int
}

func MustEncode(v any) io.Reader {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

func setupRouter(t *testing.T) http.Handler {
	t.Helper()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	conn, err := database.InitDB(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	t.Cleanup(func() { conn.Close() })

	err = database.RunMigrations(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running migrations: %v\n", err)
		os.Exit(1)
	}

	wom := utils.NewWomClient(cfg)

	srv, err := handlers.NewServer(conn, wom, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating server: %v\n", err)
		os.Exit(1)
	}

	// No auth/rate-limit middleware — just routes
	r := chi.NewRouter()
	routes.AttachV1Routes(r, srv)

	return r
}

func TestMain(t *testing.T) {
	router := setupRouter(t)

	vars := TestVariables{
		GuildId:         "test_guild",
		UserId:          "test_user",
		Rsn:             "Comfy hug",
		RsnExtra:        "Uncomfy hug",
		ChannelId:       "2012",
		Multiplier:      1,
		WomId:           "39527",
		EventClassicId:  77922,
		EventTeamId:     66321,
		AchievementName: "Ironman",
	}

	createUser := TestTable{
		Name:   "Create User",
		Method: "POST",
		Path:   fmt.Sprintf("/api/v1/guilds/%s/users", vars.GuildId),
		Body: models.CreateUserBody{
			UserId: vars.UserId,
			RSN:    vars.Rsn,
		},
		StatusCode: 200,
	}

	tt := []TestTable{
		// Guild CRUD
		{
			Name:   "Create Guild",
			Method: "POST",
			Path:   "/api/v1/guilds",
			Body: models.InputGuild{
				GuildId: vars.GuildId,
			},
			StatusCode: 200,
		},
		{
			Name:       "Guild Exists",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s", vars.GuildId),
			StatusCode: 200,
		},

		// Misc
		{
			Name:       "Get Bosses",
			Method:     "GET",
			Path:       "/api/v1/bosses",
			StatusCode: 200,
		},
		{
			Name:       "Get Categories",
			Method:     "GET",
			Path:       "/api/v1/categories",
			StatusCode: 200,
		},
		{
			Name:       "Get Achievements",
			Method:     "GET",
			Path:       "/api/v1/achievements",
			StatusCode: 200,
		},

		// User CRUD
		createUser,
		{
			Name:       "Single User Exists (By User ID)",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},
		{
			Name:       "Single User Exists (By WOM ID)",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/wom/%s", vars.GuildId, vars.WomId),
			StatusCode: 200,
		},
		{
			Name:       "Single User Exists (By RSN)",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/rsn/%s", vars.GuildId, vars.RsnEscaped()),
			StatusCode: 200,
		},
		{
			Name:       "Get User Events",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/events", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},
		{
			Name:       "Get User Times",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/times", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},
		{
			Name:       "Get User Achievements",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},

		// RSN
		{
			Name:   "Create RSN",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns", vars.GuildId, vars.UserId),
			Body: models.CreateRsnBody{
				RSN: vars.RsnExtra,
			},
			StatusCode: 200,
		},
		{
			Name:       "Delete RSN",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns/%s", vars.GuildId, vars.UserId, vars.RsnExtraEscaped()),
			StatusCode: 200,
		},

		// Points
		{
			Name:       "Update Points (Event)",
			Method:     "PUT",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/points/split_high", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},
		{
			Name:       "Update Points (Custom)",
			Method:     "PUT",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/points/custom/30", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},
		{
			Name:       "Get Point Sources",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/points", vars.GuildId),
			StatusCode: 200,
		},

		// Leaderboard
		{
			Name:       "Leaderboard Exists",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/leaderboard", vars.GuildId),
			StatusCode: 200,
		},

		// WOM
		{
			Name:       "End Competition",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/wom/competition/66321/cutoff/30", vars.GuildId),
			StatusCode: 200,
		},

		// Events
		{
			Name:   "Create Classic Event",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", vars.GuildId),
			Body: models.InputEvent{
				EventId:        vars.EventClassicId,
				PositionCutoff: 5,
			},
			StatusCode: 200,
		},
		{
			Name:   "Create Team Event",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", vars.GuildId),
			Body: models.InputEvent{
				EventId: vars.EventTeamId,
				TeamNames: []string{
					"The Jack Off Lanter",
					"Green Fingerers",
				},
			},
			StatusCode: 200,
		},
		{
			Name:       "List Events",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events", vars.GuildId),
			StatusCode: 200,
		},
		{
			Name:       "Get Event",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", vars.GuildId, vars.EventClassicId),
			StatusCode: 200,
		},
		{
			Name:       "Delete Classic Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", vars.GuildId, vars.EventClassicId),
			StatusCode: 200,
		},
		{
			Name:       "Delete Team Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", vars.GuildId, vars.EventTeamId),
			StatusCode: 200,
		},

		// Achievement give/remove
		{
			Name:       "Give Achievement",
			Method:     "POST",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements/%s", vars.GuildId, vars.UserId, vars.AchievementName),
			StatusCode: 200,
		},
		{
			Name:       "Remove Achievement",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements/%s", vars.GuildId, vars.UserId, vars.AchievementName),
			StatusCode: 200,
		},

		// Times
		{
			Name:   "Create Time",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/times", vars.GuildId),
			Body: models.InputTime{
				Time:     rand.Int(),
				BossName: "vardorvis",
				UserIds:  []string{vars.UserId},
			},
			StatusCode: 200,
		},

		// Guild times
		{
			Name:       "Get Guild Times",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/times", vars.GuildId),
			StatusCode: 200,
		},

		// Delete user variations
		{
			Name:       "Delete User (By User ID)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s", vars.GuildId, vars.UserId),
			StatusCode: 200,
		},
		createUser,
		{
			Name:       "Delete User (By Wom ID)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/wom/%s", vars.GuildId, vars.WomId),
			StatusCode: 200,
		},
		createUser,
		{
			Name:       "Delete User (By RSN)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/rsn/%s", vars.GuildId, vars.RsnEscaped()),
			StatusCode: 200,
		},

		// Cleanup
		{
			Name:       "Delete Guild",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s", vars.GuildId),
			StatusCode: 200,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()

			var body io.Reader
			if tc.Body != nil {
				body = MustEncode(tc.Body)
			}

			r := httptest.NewRequest(tc.Method, tc.Path, body)
			if tc.Body != nil {
				r.Header.Set("Content-Type", "application/json")
			}

			router.ServeHTTP(w, r)

			if w.Code != tc.StatusCode {
				t.Logf("Response: %s", w.Body.String())
				t.Fatalf("Expected status %d, got %d", tc.StatusCode, w.Code)
			}
		})
	}
}
