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

	"tectonic-api/database"
	"tectonic-api/handlers"
	"tectonic-api/models"

	"github.com/gorilla/mux"
)

type TestVariables struct {
	GuildId        string
	UserId         string
	Rsn            string
	RsnExtra       string
	ChannelId      string
	Multiplier     int
	WomId          string
	EventClassicId string
	EventTeamId    string
}

func (tv TestVariables) RsnEscaped() string {
	return url.PathEscape(tv.Rsn)
}

func (tv TestVariables) RsnExtraEscaped() string {
	return url.PathEscape(tv.Rsn)
}

type TestTable struct {
	Name       string
	Method     string
	Path       string
	Vars       map[string]string
	Body       any
	Handler    http.HandlerFunc
	StatusCode int
}

func MustEncode[T any](t T) io.Reader {
	byte, err := json.Marshal(t)
	buf := bytes.NewBuffer(byte)
	if err != nil {
		panic(err)
	}

	return buf
}

func TestMain(t *testing.T) {
	conn, err := database.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = database.RunMigrations(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running migrations: %v\n", err)
		os.Exit(1)
	}

	handlers.InitHandlers(conn)
	vars := TestVariables{
		GuildId:        "420",
		UserId:         "69",
		Rsn:            "Comfy hug",
		RsnExtra:       "Uncomfy hug",
		ChannelId:      "2012",
		Multiplier:     1,
		WomId:          "39527",
		EventClassicId: "77922",
		EventTeamId: "66321",
	}

	createUser := TestTable{
		Name:   "Create User",
		Method: "POST",
		Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s", vars.GuildId, vars.UserId),
		Vars: map[string]string{
			"guild_id": vars.GuildId,
			"user_id":  vars.UserId,
		},
		Body: models.InputUser{
			UserId:  vars.UserId,
			GuildId: vars.GuildId,
			RSN:     vars.Rsn,
		},
		Handler:    handlers.CreateUser,
		StatusCode: 201,
	}

	tt := []TestTable{
		{
			Name:   "Create Guild",
			Method: "POST",
			Path:   "/api/v1/guilds",
			Body: models.InputGuild{
				GuildId: vars.GuildId,
			},
			Handler:    handlers.CreateGuild,
			StatusCode: 201,
		},
		{
			Name:   "Update Guild",
			Method: "PUT",
			Path:   fmt.Sprintf("/api/v1/guilds/%s", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Body: models.UpdateGuild{
				GuildId:     vars.GuildId,
				Multiplier:  vars.Multiplier,
				PbChannelId: vars.ChannelId,
			},
			Handler:    handlers.UpdateGuild,
			StatusCode: 204,
		},
		{
			Name:   "Guild Exists",
			Method: "GET",
			Path:   "/api/v1/guilds",
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Handler:    handlers.GetGuild,
			StatusCode: 200,
		},
		{
			Name:       "Get Bosses",
			Method:     "GET",
			Path:       "/api/v1/bosses",
			Handler:    handlers.GetBosses,
			StatusCode: 200,
		},
		{
			Name:       "Get Categories",
			Method:     "GET",
			Path:       "/api/v1/categories",
			Handler:    handlers.GetCategories,
			StatusCode: 200,
		},
		createUser,
		{
			Name:   "Create RSN",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns/%s", vars.GuildId, vars.UserId, vars.RsnExtraEscaped()),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"user_id":  vars.UserId,
				"rsn":      vars.RsnExtra,
			},
			Handler:    handlers.CreateRSN,
			StatusCode: 204,
		},
		{
			Name:   "Delete RSN",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s/rsns/%s", vars.GuildId, vars.UserId, vars.RsnExtraEscaped()),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"user_id":  vars.UserId,
				"rsn":      vars.RsnExtra,
			},
			Handler:    handlers.RemoveRSN,
			StatusCode: 204,
		},
		{
			Name:   "Single User Exists (By User ID)",
			Method: "GET",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s", vars.GuildId, vars.UserId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"user_ids": vars.UserId,
			},
			Handler:    handlers.GetUsersById,
			StatusCode: 200,
		},
		{
			Name:   "Single User Exists (By WOM ID)",
			Method: "GET",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/wom/%s", vars.GuildId, vars.WomId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"wom_ids":  vars.WomId,
			},
			Handler:    handlers.GetUsersByWom,
			StatusCode: 200,
		},
		{
			Name:   "Single User Exists (By RSN)",
			Method: "GET",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/rsn/%s", vars.GuildId, vars.RsnEscaped()),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"rsns":     vars.Rsn,
			},
			Handler:    handlers.GetUsersByRsn,
			StatusCode: 200,
		},
		{
			Name:   "Update Points (Event)",
			Method: "PUT",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s/points/split_high", vars.GuildId, vars.UserId),
			Vars: map[string]string{
				"point_event": "split_high",
				"guild_id":    vars.GuildId,
				"user_ids":    vars.UserId,
			},
			Handler:    handlers.UpdatePoints,
			StatusCode: 200,
		},
		{
			Name:   "Update Points (Custom)",
			Method: "PUT",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s/points/custom/30", vars.GuildId, vars.UserId),
			Vars: map[string]string{
				"points":   "30",
				"guild_id": vars.GuildId,
				"user_ids": vars.UserId,
			},
			Handler:    handlers.UpdatePointsCustom,
			StatusCode: 200,
		},
		{
			Name:   "Leaderboard Exists",
			Method: "GET",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/leaderboard", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Handler:    handlers.GetLeaderboard,
			StatusCode: 200,
		},
		{
			Name:   "End Competition",
			Method: "GET",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/wom/competition/66321/cutoff/30", vars.GuildId),
			Vars: map[string]string{
				"guild_id":       vars.GuildId,
				"competition_id": "66321",
				"cutoff":         "30",
			},
			Handler:    handlers.EndCompetition,
			StatusCode: 200,
		},
		{
			Name:   "Create Classic Event",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Body: models.InputEvent{
				EventId:        vars.EventClassicId,
				PositionCutoff: 5,
			},
			Handler:    handlers.RegisterEvent,
			StatusCode: 201,
		},
		{
			Name:   "Create Team Event",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Body: models.InputEvent{
				EventId:        vars.EventTeamId,
				TeamNames: []string{
					"The Jack Off Lanter",
					"Green Fingerers",
				},
			},
			Handler:    handlers.RegisterEvent,
			StatusCode: 201,
		},
		{
			Name:   "List Events",
			Method: "GET",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Handler:    handlers.GetEvents,
			StatusCode: 200,
		},
		{
			Name:   "Delete Classic Events",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events/%s", vars.GuildId, vars.EventClassicId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"event_id": vars.EventClassicId,
			},
			Handler:    handlers.DeleteEvent,
			StatusCode: 200,
		},
		{
			Name:   "Delete Team Events",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/events/%s", vars.GuildId, vars.EventClassicId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"event_id": vars.EventTeamId,
			},
			Handler:    handlers.DeleteEvent,
			StatusCode: 200,
		},
		{
			Name:   "Create Time",
			Method: "POST",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/times", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Body: models.InputTime{
				GuildId:  vars.GuildId,
				Time:     rand.Int(),
				BossName: "vardorvis",
				UserIds:  []string{vars.UserId},
			},
			Handler:    handlers.CreateTime,
			StatusCode: 201,
		},
		{
			Name:   "Delete User (By User ID)",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/users/%s", vars.GuildId, vars.UserId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"user_id":  vars.UserId,
			},
			Handler:    handlers.RemoveUserById,
			StatusCode: 204,
		},
		createUser,
		{
			Name:   "Delete User (By Wom ID)",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/wom/%s", vars.GuildId, vars.WomId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"wom_id":   vars.WomId,
			},
			Handler:    handlers.RemoveUserByWom,
			StatusCode: 204,
		},
		createUser,
		{
			Name:   "Delete User (By RSN)",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s/rsn/%s", vars.GuildId, vars.RsnEscaped()),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
				"rsn":      vars.Rsn,
			},
			Handler:    handlers.RemoveUserByRsn,
			StatusCode: 204,
		},
		{
			Name:   "Delete Guild",
			Method: "DELETE",
			Path:   fmt.Sprintf("/api/v1/guilds/%s", vars.GuildId),
			Vars: map[string]string{
				"guild_id": vars.GuildId,
			},
			Handler:    handlers.DeleteGuild,
			StatusCode: 204,
		},
	}

	for _, exp := range tt {
		t.Run(exp.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := MustEncode(exp.Body)

			r := httptest.NewRequest(exp.Method, exp.Path, body)
			if exp.Vars != nil {
				r = mux.SetURLVars(r, exp.Vars)
			}

			exp.Handler(w, r)

			if w.Result().StatusCode != exp.StatusCode {
				t.Logf("%s", w.Body.String())
				t.Fatalf("Expected status code %d, got %d", exp.StatusCode, w.Result().StatusCode)
			}
		})
	}
}
