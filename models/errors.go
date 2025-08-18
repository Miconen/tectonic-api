package models

import (
	"net/http"
)

type APIV1Error interface {
	Message() string
	Status() int
	ToErrorResponse() ErrorResponse
}

type APIV1ErrorWithDetails interface {
	APIV1Error
	WithDetails(details any) APIV1ErrorWithDetails
}

type APIV1ErrorTodo struct {
	detail string
	status int
}

type ValidationError struct {
	Code    APIV1ErrorCode
	details any
}

func NewValidationError(details any) APIV1ErrorWithDetails {
	return &ValidationError{
		Code:    ERROR_VALIDATION_FAILED,
		details: details,
	}
}

func (e *ValidationError) Message() string {
	return ERROR_VALIDATION_FAILED.String()
}

func (e *ValidationError) Status() int {
	return http.StatusBadRequest
}

func (e *ValidationError) ToErrorResponse() ErrorResponse {
	return ErrorResponse{
		Code:    uint(e.Code),
		Message: e.Message(),
		Details: e.details,
	}
}

func (e *ValidationError) WithDetails(details any) APIV1ErrorWithDetails {
	return &ValidationError{
		Code:    e.Code,
		details: details,
	}
}

func ERROR_TODO(status int, detail string) APIV1Error {
	return &APIV1ErrorTodo{status: status, detail: detail}
}

func (et *APIV1ErrorTodo) Message() string {
	return untreated.String()
}

func (et *APIV1ErrorTodo) Status() int {
	return et.status
}

func (et *APIV1ErrorTodo) ToErrorResponse() ErrorResponse {
	return ErrorResponse{
		Code:    uint(untreated),
		Message: et.Message(),
		Details: et.detail,
	}
}

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

	ERROR_TEAM_NOT_FOUND // User haven't participated in the run
	ERROR_TEAM_EXISTS    // User already associated to this run

	ERROR_EVENT_NOT_FOUND // Event not found
	ERROR_EVENT_EXISTS    // Event already exists

	ERROR_ACHIEVEMENT_NOT_FOUND // Achievement not found
	ERROR_ACHIEVEMENT_EXISTS    // Achievement already exists

	ERROR_USER_ACHIEVEMENT_NOT_FOUND // User doesn't have this achievement
	ERROR_USER_ACHIEVEMENT_EXISTS    // User already has this achievement

	ERROR_POINT_SOURCE_NOT_FOUND // Point source not found
	ERROR_POINT_SOURCE_EXISTS    // Point source already exists
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
	case ERROR_API_UNAVAILABLE:
		return http.StatusInternalServerError
	case ERROR_API_DEAD:
		return http.StatusInternalServerError
	case ERROR_WRONG_BODY:
		return http.StatusBadRequest
	case ERROR_WRONG_PARAMS:
		return http.StatusBadRequest
	case ERROR_VALIDATION_FAILED:
		return http.StatusBadRequest
	// Model
	case ERROR_GUILD_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_GUILD_EXISTS:
		return http.StatusConflict

	case ERROR_BOSS_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_BOSS_EXISTS:
		return http.StatusConflict

	case ERROR_CATEGORY_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_CATEGORY_EXISTS:
		return http.StatusConflict

	case ERROR_RSN_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_RSN_EXISTS:
		return http.StatusConflict

	case ERROR_TIME_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_TIME_EXISTS:
		return http.StatusConflict

	case ERROR_USER_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_USER_EXISTS:
		return http.StatusConflict

	case ERROR_WOMID_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_WOMID_EXISTS:
		return http.StatusConflict

	case ERROR_PARTICIPATION_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_PARTICIPATION_EXISTS:
		return http.StatusConflict

	case ERROR_GUILD_BOSS_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_GUILD_BOSS_EXISTS:
		return http.StatusConflict

	case ERROR_GUILD_CATEGORY_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_GUILD_CATEGORY_EXISTS:
		return http.StatusConflict

	case ERROR_TEAM_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_TEAM_EXISTS:
		return http.StatusConflict

	case ERROR_EVENT_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_EVENT_EXISTS:
		return http.StatusConflict

	case ERROR_ACHIEVEMENT_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_ACHIEVEMENT_EXISTS:
		return http.StatusConflict

	case ERROR_USER_ACHIEVEMENT_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_USER_ACHIEVEMENT_EXISTS:
		return http.StatusConflict

	case ERROR_POINT_SOURCE_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_POINT_SOURCE_EXISTS:
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
