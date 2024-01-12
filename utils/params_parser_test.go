package utils

import (
	"net/http"
	"net/url"
	"testing"
)

func TestParseParametersSuccess(t *testing.T) {
	r := &http.Request{
		URL: &url.URL{
			RawQuery: "user_id=10&guild_id=9",
		},
	}

	vals, err := ParseParametersURL(r, "user_id", "guild_id")
	if err != nil {
		t.Fail()
		t.Log(err.Error())
	}

	if vals["user_id"] != "10" {
		t.Errorf("user_id = %s; want 10", vals["user_id"])
	}
	if vals["guild_id"] != "9" {
		t.Errorf("guild_id = %s; want 9", vals["guild_id"])
	}
}

func TestParseParametersFail(t *testing.T) {
	r := []http.Request{
		{
			URL: &url.URL{
				RawQuery: "user_id=10&guild_id=",
			},
		},
		{
			URL: &url.URL{
				RawQuery: "user_id=10&guild_id=9&guild_id=12",
			},
		},
	}

	for i := range r {
		_, err := ParseParametersURL(&r[i], "user_id", "guild_id")
		if err == nil {
			t.Fail()
		}
	}
}

func TestParseParametersOptional(t *testing.T) {
	r := &http.Request{
		URL: &url.URL{
			RawQuery: "user_id=10&guild_id=9&meaning_of_life=42",
		},
	}

	vals, err := ParseParametersURL(r, "user_id", "guild_id")
	if err != nil {
		t.Fail()
		t.Log(err.Error())
	}

	if val, ok := vals["meaning_of_life"]; !ok || val != "42" {
		t.Fail()
		t.Logf("meaning_of_life = %s; want 42", val)
	}
}
