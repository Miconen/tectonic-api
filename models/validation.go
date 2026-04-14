// models/validation.go
package models

import (
	"github.com/danielgtaylor/huma/v2"
)

// Discord Snowflake — 17-19 digit numeric string, no leading zero
type DiscordSnowflake string

func (d DiscordSnowflake) String() string { return string(d) }

func (d DiscordSnowflake) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	s := string(d)
	if len(s) < 17 || len(s) > 19 {
		return []error{&huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "must be a valid Discord ID (17-19 digits)",
			Value:    s,
		}}
	}
	if s[0] == '0' {
		return []error{&huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "Discord ID cannot start with zero",
			Value:    s,
		}}
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return []error{&huma.ErrorDetail{
				Location: prefix.String(),
				Message:  "Discord ID must be numeric",
				Value:    s,
			}}
		}
	}
	return nil
}

// RSN — 1-12 chars, letters/numbers/spaces/hyphens/underscores
type RSN string

func (r RSN) String() string { return string(r) }

func (r RSN) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	s := string(r)
	if len(s) == 0 || len(s) > 12 {
		return []error{&huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "RSN must be 1-12 characters",
			Value:    s,
		}}
	}
	if s[0] == ' ' || s[len(s)-1] == ' ' {
		return []error{&huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "RSN cannot start or end with a space",
			Value:    s,
		}}
	}
	for _, c := range s {
		if !isValidRSNChar(c) {
			return []error{&huma.ErrorDetail{
				Location: prefix.String(),
				Message:  "RSN contains invalid characters (allowed: letters, numbers, spaces, hyphens, underscores)",
				Value:    s,
			}}
		}
	}
	return nil
}

func isValidRSNChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == ' ' || r == '-' || r == '_'
}

func SnowflakesToStrings(ids []DiscordSnowflake) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = string(id)
	}
	return out
}
