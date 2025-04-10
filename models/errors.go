package models

import "net/http"

type APIV1Error uint

// Request-based errors
const (
	ERROR_WRONG_PARAMS APIV1Error = iota
	ERROR_WRONG_BODY
	ERROR_INVALID_TOKEN
	ERROR_DATA_DOESNT_EXISTS
)

// Model-based errors
const (
	ERROR_GUILD_NOT_FOUND APIV1Error = 1000 + iota
	ERROR_GUILD_EXISTS

	ERROR_BOSS_NOT_FOUND
	ERROR_CATEGORY_NOT_FOUND

	ERROR_RSN_NOT_FOUND
	ERROR_RSN_EXISTS

	ERROR_TIME_NOT_FOUND

	ERROR_USER_NOT_FOUND
	ERROR_USER_EXISTS

	ERROR_WOMID_NOT_FOUND
	ERROR_WOMID_EXISTS

	ERROR_PARTICIPATION_NOT_FOUND

	ERROR_EVENT_NOT_FOUND
	ERROR_EVENT_EXISTS

	ERROR_ACHIEVEMENT_NOT_FOUND
	ERROR_ACHIEVEMENT_EXISTS
)

// Server errors
const (
	ERROR_API_UNAVAILABLE APIV1Error = 2000 + iota
	ERROR_API_DEAD
	ERROR_WOM_UNAVAILABLE
)

func (e APIV1Error) Message() (message string) {
	switch e {
	case ERROR_INVALID_TOKEN:
		message = "Your token is invalid"
	case ERROR_BOSS_NOT_FOUND:
		message = "Boss have not been found"
	case ERROR_CATEGORY_NOT_FOUND:
		message = "Category have not been found"
	case ERROR_GUILD_EXISTS:
		message = "Guild specified already exists"
	case ERROR_GUILD_NOT_FOUND:
		message = "Guild have not been found"
	case ERROR_RSN_EXISTS:
		message = "RSN specified is taken by another user"
	case ERROR_RSN_NOT_FOUND:
		message = "RSN have not been found"
	case ERROR_TIME_NOT_FOUND:
		message = "Time have not been found"
	case ERROR_USER_EXISTS:
		message = "User specified already exists"
	case ERROR_USER_NOT_FOUND:
		message = "User have not been found"
	case ERROR_WOMID_EXISTS:
		message = "Wise old man Id specified already exists"
	case ERROR_WOMID_NOT_FOUND:
		message = "Wise old man Id have not been found"
	case ERROR_WOM_UNAVAILABLE:
		message = "Wise old man is unavaliable, please try again later"
	case ERROR_API_UNAVAILABLE:
		message = "Something's wrong with the API, please try again later"
	case ERROR_API_DEAD:
		message = "Server broke, contact us if you see this error"
	case ERROR_WRONG_BODY:
		message = "Body is malformated, please check docs for example on how to send the request"
	case ERROR_WRONG_PARAMS:
		message = "Params are malformated, please check docs for example on how to send the request"
	case ERROR_PARTICIPATION_NOT_FOUND:
		message = "No participations found with specified cutoff"
	case ERROR_DATA_DOESNT_EXISTS:
		message = "Data sent seems correct, but couldn't find any correspondence in our end. Please, verify its veracity on their GET routes respectively"
	case ERROR_EVENT_NOT_FOUND:
		message = "Event have not been found"
	case ERROR_EVENT_EXISTS:
		message = "Event specified already exists"
	case ERROR_ACHIEVEMENT_NOT_FOUND:
		message = "Achievement have not been found"
	case ERROR_ACHIEVEMENT_EXISTS:
		message = "Achievement specified already exists"
	}

	return
}

func (e APIV1Error) Status() int {
	switch e {
	case ERROR_INVALID_TOKEN:
		return http.StatusUnauthorized
	case ERROR_BOSS_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_CATEGORY_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_GUILD_EXISTS:
		return http.StatusConflict
	case ERROR_GUILD_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_EVENT_EXISTS:
		return http.StatusConflict
	case ERROR_EVENT_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_ACHIEVEMENT_EXISTS:
		return http.StatusConflict
	case ERROR_ACHIEVEMENT_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_RSN_EXISTS:
		return http.StatusConflict
	case ERROR_RSN_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_TIME_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_USER_EXISTS:
		return http.StatusConflict
	case ERROR_USER_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_WOMID_EXISTS:
		return http.StatusConflict
	case ERROR_WOMID_NOT_FOUND:
		return http.StatusNotFound
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
	case ERROR_PARTICIPATION_NOT_FOUND:
		return http.StatusNotFound
	case ERROR_DATA_DOESNT_EXISTS:
		return http.StatusNotFound
	default:
		return http.StatusOK
	}
}

func (e APIV1Error) ToErrorResponse() ErrorResponse {
	msg := e.Message()
	return ErrorResponse{
		Code:    uint(e),
		Message: msg,
	}
}
