package models

import (
	"net/http"
)

type APIV1Error interface {
	Message() string
	Status() int
	Code() uint
}

func (e APIV1ErrorCode) Code() uint      { return uint(e) }
func (e APIV1ErrorCode) Message() string { return e.String() }

//go:generate stringer -type=APIV1ErrorCode -linecomment -output=errors_string.go
type APIV1ErrorCode uint

// Request-based errors
const (
	ERROR_WRONG_PARAMS      APIV1ErrorCode = iota // Params are malformated, please check docs for example on how to send the request
	ERROR_WRONG_BODY                              // Body is malformated, please check docs for example on how to send the request
	ERROR_VALIDATION_FAILED                       // Request validation failed
	ERROR_INVALID_TOKEN                           // Your token is invalid
)

// Model-based errors
const (
	ERROR_GUILD_NOT_FOUND APIV1ErrorCode = 1000 + iota // Guild not found
	ERROR_GUILD_EXISTS                                 // Guild already exists

	ERROR_BOSS_NOT_FOUND // Boss not found
	ERROR_BOSS_EXISTS    // Boss already exists

	ERROR_CATEGORY_NOT_FOUND // Category not found
	ERROR_CATEGORY_EXISTS    // Category already exists

	ERROR_RSN_NOT_FOUND // RSN not found
	ERROR_RSN_EXISTS    // RSN is already associated to a user

	ERROR_TIME_NOT_FOUND // Time not found
	ERROR_TIME_EXISTS    // Time already exists

	ERROR_USER_NOT_FOUND // User not found
	ERROR_USER_EXISTS    // User already exists

	ERROR_WOMID_NOT_FOUND // Wom Id not found
	ERROR_WOMID_EXISTS    // Wom Id already associated do an RSN

	ERROR_PARTICIPATION_NOT_FOUND // User not found in event
	ERROR_PARTICIPATION_EXISTS    // User already registered in the event

	ERROR_GUILD_BOSS_NOT_FOUND // Guild doesn't have this boss registered
	ERROR_GUILD_BOSS_EXISTS    // Guild already have this boss registered

	ERROR_GUILD_CATEGORY_NOT_FOUND // Guild doesn't have this category registered
	ERROR_GUILD_CATEGORY_EXISTS    // Guild already have this category registered

	ERROR_TEAM_NOT_FOUND // User has not participated in the run
	ERROR_TEAM_EXISTS    // User already associated to this run

	ERROR_EVENT_NOT_FOUND // Event not found
	ERROR_EVENT_EXISTS    // Event already exists

	ERROR_ACHIEVEMENT_NOT_FOUND // Achievement not found
	ERROR_ACHIEVEMENT_EXISTS    // Achievement already exists

	ERROR_USER_ACHIEVEMENT_NOT_FOUND // User doesn't have this achievement
	ERROR_USER_ACHIEVEMENT_EXISTS    // User already has this achievement

	ERROR_POINT_SOURCE_NOT_FOUND // Point source not found
	ERROR_POINT_SOURCE_EXISTS    // Point source already exists

	ERROR_COMBAT_ACHIEVEMENT_NOT_FOUND // Combat achievement not found
	ERROR_COMBAT_ACHIEVEMENT_EXISTS    // Combat achievement already exists
)

// Server errors
const (
	ERROR_API_UNAVAILABLE APIV1ErrorCode = 2000 + iota // Something's wrong with the API, please try again later
	ERROR_API_DEAD                                     // Server broke, contact us if you see this error
	ERROR_WOM_UNAVAILABLE                              // Wise old man is unavaliable, please try again later
)

const untreated APIV1ErrorCode = 10000 // Some error happened but left untreated, please file an issue here: https://github.com/Miconen/tectonic-api/issues/new

func (e APIV1ErrorCode) Message() string {
	return e.String()
}

func (e APIV1ErrorCode) Status() int {
	switch e {
	case ERROR_INVALID_TOKEN:
		return http.StatusUnauthorized
	case ERROR_WOM_UNAVAILABLE:
		return http.StatusServiceUnavailable
	case ERROR_API_UNAVAILABLE, ERROR_API_DEAD:
		return http.StatusInternalServerError
	case ERROR_WRONG_BODY, ERROR_WRONG_PARAMS, ERROR_VALIDATION_FAILED:
		return http.StatusBadRequest
	}

	switch e {
	case ERROR_GUILD_NOT_FOUND,
		ERROR_BOSS_NOT_FOUND,
		ERROR_CATEGORY_NOT_FOUND,
		ERROR_RSN_NOT_FOUND,
		ERROR_TIME_NOT_FOUND, ERROR_USER_NOT_FOUND,
		ERROR_WOMID_NOT_FOUND,
		ERROR_PARTICIPATION_NOT_FOUND,
		ERROR_GUILD_BOSS_NOT_FOUND,
		ERROR_GUILD_CATEGORY_NOT_FOUND,
		ERROR_TEAM_NOT_FOUND,
		ERROR_EVENT_NOT_FOUND,
		ERROR_ACHIEVEMENT_NOT_FOUND,
		ERROR_USER_ACHIEVEMENT_NOT_FOUND,
		ERROR_POINT_SOURCE_NOT_FOUND,
		ERROR_COMBAT_ACHIEVEMENT_NOT_FOUND:
		return http.StatusNotFound

	case ERROR_GUILD_EXISTS,
		ERROR_BOSS_EXISTS,
		ERROR_CATEGORY_EXISTS,
		ERROR_RSN_EXISTS,
		ERROR_TIME_EXISTS,
		ERROR_USER_EXISTS,
		ERROR_WOMID_EXISTS,
		ERROR_PARTICIPATION_EXISTS,
		ERROR_GUILD_BOSS_EXISTS,
		ERROR_GUILD_CATEGORY_EXISTS,
		ERROR_TEAM_EXISTS,
		ERROR_EVENT_EXISTS,
		ERROR_ACHIEVEMENT_EXISTS,
		ERROR_USER_ACHIEVEMENT_EXISTS,
		ERROR_POINT_SOURCE_EXISTS,
		ERROR_COMBAT_ACHIEVEMENT_EXISTS:
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}

func (e APIV1ErrorCode) ToErrorResponse() ErrorResponse {
	msg := e.Message()
	return ErrorResponse{
		Code:    uint(e),
		Message: msg,
	}
}

func ValidationFailed(details any) APIV1Error {
	return NewValidationError(details)
}
