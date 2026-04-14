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
	GuildID         string
	UserID          string
	Rsn             string
	RsnExtra        string
	ChannelID       string
	WomID           string
	EventClassicID  int
	EventTeamID     int
	AchievementName string
}

func (tv TestVariables) RsnEscaped() string {
	return url.PathEscape(tv.Rsn)
}

func (tv TestVariables) RsnExtraEscaped() string {
	return url.PathEscape(tv.RsnExtra)
}

type TestCase struct {
	Name       string
	Method     string
	Path       string
	Body       any
	StatusCode int
}

func mustEncode(v any) io.Reader {
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

	r := chi.NewRouter()
	routes.AttachV1Routes(r, srv)
	return r
}

func TestRoutes(t *testing.T) {
	router := setupRouter(t)

	v := TestVariables{
		GuildID:         "123456789012345678",
		UserID:          "987654321098765432",
		Rsn:             "Comfy hug",
		RsnExtra:        "Uncomfy hug",
		ChannelID:       "111222333444555666",
		WomID:           "39527",
		EventClassicID:  77922,
		EventTeamID:     66321,
		AchievementName: "Ironman",
	}

	createUser := TestCase{
		Name:   "Create User",
		Method: "POST",
		Path:   fmt.Sprintf("/api/v1/guilds/%s/users", v.GuildID),
		Body: models.CreateUserBody{
			UserID: models.DiscordSnowflake(v.UserID),
			RSN:    models.RSN(v.Rsn),
		},
		StatusCode: 200,
	}

	tt := []TestCase{
		// === Guild ===
		{
			Name:   "Create Guild",
			Method: "POST",
			Path:   "/api/v1/guilds",
			Body: models.InputGuild{
				GuildID: models.DiscordSnowflake(v.GuildID),
			},
			StatusCode: 200,
		},
		{
			Name:       "Get Guild",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s", v.GuildID),
			StatusCode: 200,
		},
		{
			Name:   "Update Guild",
			Method: "PUT",
			Path:   fmt.Sprintf("/api/v1/guilds/%s", v.GuildID),
			Body: models.UpdateGuildBody{
				ModChannelID: ptrTo(models.DiscordSnowflake(v.ChannelID)),
			},
			StatusCode: 200,
		},

		// === Misc ===
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

		// === User ===
		createUser,
		{
			Name:       "Get User (By ID)",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s", v.GuildID, v.UserID),
			StatusCode: 200,
		},
		{
			Name:       "Get User (By WOM)",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/wom/%s", v.GuildID, v.WomID),
			StatusCode: 200,
		},
		{
			Name:       "Get User (By RSN)",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/rsn/%s", v.GuildID, v.RsnEscaped()),
			StatusCode: 200,
		},
		{
			Name:       "Get User Events",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/events", v.GuildID, v.UserID),
			StatusCode: 200,
		},
		{
			Name:       "Get User Times",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/times", v.GuildID, v.UserID),
			StatusCode: 200,
		},
		{
			Name:       "Get User Achievements",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements", v.GuildID, v.UserID),
			StatusCode: 200,
		},

		// === RSN ===
		{
			Name:   "Create RSN",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns", v.GuildID, v.UserID),
			Body: models.CreateRsnBody{
				RSN: models.RSN(v.RsnExtra),
			},
			StatusCode: 200,
		},
		{
			Name:       "Delete RSN",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns/%s", v.GuildID, v.UserID, v.RsnExtraEscaped()),
			StatusCode: 200,
		},

		// === Points ===
		{
			Name:       "Update Points (Event)",
			Method:     "PUT",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/points/split_high", v.GuildID, v.UserID),
			StatusCode: 200,
		},
		{
			Name:       "Update Points (Custom)",
			Method:     "PUT",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/points/custom/30", v.GuildID, v.UserID),
			StatusCode: 200,
		},
		{
			Name:       "Get Point Sources",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/points", v.GuildID),
			StatusCode: 200,
		},

		// === Leaderboard ===
		{
			Name:       "Get Leaderboard",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/leaderboard", v.GuildID),
			StatusCode: 200,
		},

		// === WOM ===
		{
			Name:       "End Competition",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/wom/competition/66321/cutoff/30", v.GuildID),
			StatusCode: 200,
		},

		// === Events ===
		{
			Name:   "Create Classic Event",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", v.GuildID),
			Body: models.InputEvent{
				EventID:        v.EventClassicID,
				PositionCutoff: 5,
			},
			StatusCode: 200,
		},
		{
			Name:   "Create Team Event",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", v.GuildID),
			Body: models.InputEvent{
				EventID: v.EventTeamID,
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
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events", v.GuildID),
			StatusCode: 200,
		},
		{
			Name:       "Get Event",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", v.GuildID, v.EventClassicID),
			StatusCode: 200,
		},
		{
			Name:       "Delete Classic Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", v.GuildID, v.EventClassicID),
			StatusCode: 200,
		},
		{
			Name:       "Delete Team Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", v.GuildID, v.EventTeamID),
			StatusCode: 200,
		},

		// === Achievements ===
		{
			Name:       "Give Achievement",
			Method:     "POST",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements/%s", v.GuildID, v.UserID, v.AchievementName),
			StatusCode: 200,
		},
		{
			Name:       "Remove Achievement",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements/%s", v.GuildID, v.UserID, v.AchievementName),
			StatusCode: 200,
		},

		// === Times ===
		{
			Name:   "Create Time",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/times", v.GuildID),
			Body: models.InputTime{
				Time:     rand.Intn(100000) + 1,
				BossName: "vardorvis",
				UserIDs:  []models.DiscordSnowflake{models.DiscordSnowflake(v.UserID)},
			},
			StatusCode: 200,
		},
		{
			Name:       "Get Guild Times",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/times", v.GuildID),
			StatusCode: 200,
		},

		// === Delete Users (all variations) ===
		{
			Name:       "Delete User (By ID)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s", v.GuildID, v.UserID),
			StatusCode: 200,
		},
		createUser,
		{
			Name:       "Delete User (By WOM)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/wom/%s", v.GuildID, v.WomID),
			StatusCode: 200,
		},
		createUser,
		{
			Name:       "Delete User (By RSN)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/rsn/%s", v.GuildID, v.RsnEscaped()),
			StatusCode: 200,
		},

		// === Cleanup ===
		{
			Name:       "Delete Guild",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s", v.GuildID),
			StatusCode: 200,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var body io.Reader
			if tc.Body != nil {
				body = mustEncode(tc.Body)
			}

			r := httptest.NewRequest(tc.Method, tc.Path, body)
			if tc.Body != nil {
				r.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			if w.Code != tc.StatusCode {
				t.Logf("Response: %s", w.Body.String())
				t.Fatalf("Expected status %d, got %d", tc.StatusCode, w.Code)
			}
		})
	}
}

// Helper for pointer fields
func ptrTo[T any](v T) *T {
	return &v
}
