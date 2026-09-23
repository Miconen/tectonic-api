package utils

import (
	"testing"

	"tectonic-api/models"
)

func TestRestoreSeparators(t *testing.T) {
	cases := []struct{ display, input, want string }{
		{"Foo Bar", "foo_bar", "Foo_Bar"},
		{"Foo Bar", "foo-bar", "Foo-Bar"},
		{"Foo Bar", "foo bar", "Foo Bar"},
		{"Foo Bar", "foobar", "Foo Bar"},
		{"Foo Bar", "FOO_BAR", "Foo_Bar"},
	}
	for _, c := range cases {
		if got := RestoreSeparators(c.display, models.RSN(c.input)); got != c.want {
			t.Errorf("%q,%q: got %q want %q", c.display, c.input, got, c.want)
		}
	}
}
