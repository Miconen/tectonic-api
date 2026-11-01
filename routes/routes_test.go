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
	// Never call the real Wise Old Man: it makes the test depend on the
	// network and on WOM's data staying the same.
	cfg.WOM.BaseURL = fakeWom(t)

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
		AchievementName: "Maxed",
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
			StatusCode: http.StatusNoContent,
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
			StatusCode: http.StatusNoContent,
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
			Name:       "Get User Records",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/records", v.GuildID, v.UserID),
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
			StatusCode: http.StatusNoContent,
		},
		{
			Name:       "Delete RSN",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns/%s", v.GuildID, v.UserID, v.RsnExtraEscaped()),
			StatusCode: http.StatusNoContent,
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
			Method:     "POST",
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
				PositionCutoff: 3,
			},
			StatusCode: http.StatusNoContent,
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
			StatusCode: http.StatusNoContent,
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
			// Fewer participants than the cutoff used to panic.
			Name:       "Create Classic Event (fewer players than cutoff)",
			Method:     "POST",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events", v.GuildID),
			Body:       models.InputEvent{EventID: 11111, PositionCutoff: 3},
			StatusCode: http.StatusNoContent,
		},
		{
			Name:       "Delete Small Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/11111", v.GuildID),
			StatusCode: http.StatusNoContent,
		},
		{
			Name:       "Delete Classic Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", v.GuildID, v.EventClassicID),
			StatusCode: http.StatusNoContent,
		},
		{
			Name:       "Delete Team Event",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/events/%d", v.GuildID, v.EventTeamID),
			StatusCode: http.StatusNoContent,
		},

		// === Achievements ===
		{
			Name:       "Give Achievement",
			Method:     "POST",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements/%s", v.GuildID, v.UserID, v.AchievementName),
			StatusCode: http.StatusNoContent,
		},
		{
			Name:       "Remove Achievement",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s/achievements/%s", v.GuildID, v.UserID, v.AchievementName),
			StatusCode: http.StatusNoContent,
		},

		// === Records ===
		{
			Name:   "Create Record",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/records", v.GuildID),
			Body: models.InputRecord{
				Value:    rand.Intn(100000) + 1,
				BossName: "vardorvis",
				UserIDs:  []models.DiscordSnowflake{models.DiscordSnowflake(v.UserID)},
			},
			StatusCode: 200,
		},
		{
			Name:       "Get Guild Records",
			Method:     "GET",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/records", v.GuildID),
			StatusCode: 200,
		},

		// === Delete Users (all variations) ===
		{
			Name:       "Delete User (By ID)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/%s", v.GuildID, v.UserID),
			StatusCode: http.StatusNoContent,
		},
		createUser,
		{
			Name:       "Delete User (By WOM)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/wom/%s", v.GuildID, v.WomID),
			StatusCode: http.StatusNoContent,
		},
		createUser,
		{
			Name:       "Delete User (By RSN)",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s/users/rsn/%s", v.GuildID, v.RsnEscaped()),
			StatusCode: http.StatusNoContent,
		},

		// === Cleanup ===
		{
			Name:       "Delete Guild",
			Method:     "DELETE",
			Path:       fmt.Sprintf("/api/v1/guilds/%s", v.GuildID),
			StatusCode: http.StatusNoContent,
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

// fakeWom serves the players and competitions the test uses, in the shape of
// the Wise Old Man v2 API.
func fakeWom(t *testing.T) string {
	t.Helper()
	players := map[string]int{"Comfy hug": 39527, "Uncomfy hug": 39528, "Third": 39529}
	participation := func(name, team string, gained float64) models.Participations {
		return models.Participations{
			PlayerID: players[name],
			TeamName: team,
			Player:   models.Player{ID: players[name], DisplayName: name},
			Progress: models.Progress{Gained: gained},
		}
	}
	competitions := map[string]models.WomCompetition{
		"77922": {ID: 77922, Title: "Classic event", Type: "classic", ParticipantCount: 3, Participations: []models.Participations{
			participation("Comfy hug", "", 100), participation("Uncomfy hug", "", 50), participation("Third", "", 10),
		}},
		"11111": {ID: 11111, Title: "Small event", Type: "classic", ParticipantCount: 2, Participations: []models.Participations{
			participation("Comfy hug", "", 100), participation("Uncomfy hug", "", 50),
		}},
		"66321": {ID: 66321, Title: "Team event", Type: "team", ParticipantCount: 3, Participations: []models.Participations{
			participation("Comfy hug", "The Jack Off Lanter", 100),
			participation("Uncomfy hug", "Green Fingerers", 50),
			participation("Third", "Green Fingerers", 10),
		}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /players/{rsn}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("rsn")
		id, ok := players[name]
		if !ok {
			http.Error(w, `{"message":"Player not found."}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": id, "displayName": name})
	})
	mux.HandleFunc("GET /competitions/{id}", func(w http.ResponseWriter, r *http.Request) {
		c, ok := competitions[r.PathValue("id")]
		if !ok {
			http.Error(w, `{"message":"Competition not found."}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(c)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv.URL
}

// Helper for pointer fields
func ptrTo[T any](v T) *T {
	return &v
}
