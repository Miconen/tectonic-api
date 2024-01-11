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
	r := &http.Request{
		URL: &url.URL{
			RawQuery: "user_id=10&guild_id=",
		},
	}

	_, err := ParseParametersURL(r, "user_id", "guild_id")
	if err == nil {
		t.Fail()
		t.Log(err.Error())
	}
}
