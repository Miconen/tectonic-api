package utils

import (
	"bytes"
	"io"
	"net/http/httptest"
	"tectonic-api/models"
	"testing"
)

func TestParseRequestBodySuccess(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/v1/guilds", io.NopCloser(bytes.NewReader([]byte(`{"guild_id":"someID","multiplier":1,"pb_channel_id":"someID"}`))))
	w := httptest.NewRecorder()

	p := models.InputGuild{}
	err := ParseRequestBody(w, r, &p)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
	}

	if p.GuildId != "someID" {
		t.Errorf("guild_id = %s; want 'someID'", p.GuildId)
	}
	if p.Multiplier != 1 {
		t.Errorf("multiplier = %d; want 1", p.Multiplier)
	}
}

func TestParseRequestBodySuccessMissingOptionalField(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/v1/guilds", io.NopCloser(bytes.NewReader([]byte(`{"guild_id":"someID"}`))))
	w := httptest.NewRecorder()

	p := models.InputGuild{
		Multiplier: 1,
	}
	err := ParseRequestBody(w, r, &p)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
	}

	if p.GuildId != "someID" {
		t.Errorf("guild_id = %s; want 'someID'", p.GuildId)
	}
	if p.Multiplier != 1 {
		t.Errorf("multiplier = %d; want 1", p.Multiplier)
	}
}

type TestInputGuild struct {
	GuildId    string `json:"guild_id"`
	Multiplier int    `json:"multiplier"`
}

func TestParseRequestBodySuccessRequiredFieldZero(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/v1/guilds", io.NopCloser(bytes.NewReader([]byte(`{"guild_id":"someID","multiplier":0}`))))
	w := httptest.NewRecorder()

	p := models.InputGuild{}
	err := ParseRequestBody(w, r, &p)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
	}

	if p.GuildId != "someID" {
		t.Errorf("guild_id = %s; want 'someID'", p.GuildId)
	}
	if p.Multiplier != 0 {
		t.Errorf("multiplier = %d; want 0", p.Multiplier)
	}
}
func TestParseRequestBodyFailExtraField(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/v1/guilds", io.NopCloser(bytes.NewReader([]byte(`{"guild_id":"someID","random_field":"hello"}`))))
	w := httptest.NewRecorder()

	p := models.InputGuild{}
	err := ParseRequestBody(w, r, &p)
	if err == nil {
		t.Fail()
		t.Log("Accepted errant field")
	}
}

func TestParseRequestBodyFailInvalidType(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/v1/guilds", io.NopCloser(bytes.NewReader([]byte(`{"guild_id":1}`))))
	w := httptest.NewRecorder()

	p := models.InputGuild{}
	err := ParseRequestBody(w, r, &p)
	if err == nil {
		t.Fail()
		t.Log("Didn't error on invalid type")
	}
}

func TestParseRequestBodyFailMissingRequiredField(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/v1/guilds", io.NopCloser(bytes.NewReader([]byte(`{"multiplier":1}`))))
	w := httptest.NewRecorder()

	p := models.InputGuild{}
	err := ParseRequestBody(w, r, &p)
	if err == nil {
		t.Fail()
		t.Log("Didn't error on missing required field")
	}
}
