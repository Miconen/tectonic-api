package utils

import (
	"testing"
)

func TestRequireOneSuccess(t *testing.T) {
	m := map[string]string{
		"guild_id": "foo",
		"wom":      "bar",
	}

	key, err := RequireOne(m, "user_id", "rsn", "wom")
	if err != nil {
		t.Fail()
		t.Log(err.Error())
	}

	if key != "wom" {
		t.Errorf("Key = %s, Value = %s; want user_id, bar", key, m[key])
	}
}

func TestRequireOneFail(t *testing.T) {
	m := map[string]string{
		"guild_id": "foo",
	}

	key, err := RequireOne(m, "user_id", "rsn", "wom")
	if key != "" {
		t.Errorf("key = %s; want empty string", key)
	}
	if err == nil {
		t.Errorf("Expected error; got none")
	}
}

func TestRequireOneAmount(t *testing.T) {
	m := map[string]string{
		"guild_id": "foo",
		"user_id":  "bar",
		"rsn":      "baz",
		"wom":      "qux",
	}

	key, err := RequireOne(m, "user_id", "rsn", "wom")
	if key != "" {
		t.Errorf("key = %s; want empty string", key)
	}
	if err == nil {
		t.Errorf("Expected error; got none")
	}
}
