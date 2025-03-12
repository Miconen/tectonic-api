package utils

type ApiError uint

const (
	ERROR_RSN_NOT_FOUND ApiError = 1000 +iota
	ERROR_USER_NOT_FOUND
	ERROR_GUILD_NOT_FOUND
)

func (e ApiError) String() string {
	switch e {
	case ERROR_GUILD_NOT_FOUND:
		return ""
	case ERROR_RSN_NOT_FOUND:
		return ""
	case ERROR_USER_NOT_FOUND:
		return ""
	}

	return ""
}
