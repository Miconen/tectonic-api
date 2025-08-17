package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"tectonic-api/logging"
	"tectonic-api/models"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	// Cache for boss names to avoid repeated DB queries during validation
	// TODO: Fetch bosses from database and validate against those
	validBosses map[string]bool
)

func init() {
	validate = validator.New()

	// Use JSON field names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register your custom validators
	validate.RegisterValidation("discord_snowflake", validateDiscordSnowflake)
	validate.RegisterValidation("rsn", validateRSN)
	validate.RegisterValidation("positive_time", validatePositiveTime)
	validate.RegisterValidation("boss_name", validateBossName)

	loadValidBosses()
}

// Validation error details structure
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Value   any    `json:"value"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func convertValidationErrors(err error) []ValidationErrorDetail {
	var errors []ValidationErrorDetail

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			errors = append(errors, ValidationErrorDetail{
				Field:   fe.Field(),
				Value:   fe.Value(),
				Tag:     fe.Tag(),
				Message: getValidationMessage(fe),
			})
		}
	}

	return errors
}

// Automatically handle validation errors and return detailed responses
// for malformed requests and invalid parameters.
//
// Validation is based on struct tags of the 3rd parameter.
func ParseAndValidateRequestBody(w http.ResponseWriter, r *http.Request, v any) error {
	if err := ParseRequestBody(w, r, v); err != nil {
		logging.Get().Info("error parsing request body", "error", err)
		jw := NewJsonWriter(w, r, http.StatusBadRequest)
		jw.WriteError(models.ERROR_WRONG_BODY)
		return err
	}

	if err := validate.Struct(v); err != nil {
		validationDetails := convertValidationErrors(err)
		validationError := models.ValidationFailed(validationDetails)

		jw := NewJsonWriter(w, r, validationError.Status())
		jw.WriteError(validationError)

		return fmt.Errorf("validation failed")
	}

	return nil
}

// Custom validators
// Discord Snowflake validator (17-19 digits)
func validateDiscordSnowflake(fl validator.FieldLevel) bool {
	snowflake := fl.Field().String()

	// Check length
	if len(snowflake) < 17 || len(snowflake) > 19 {
		return false
	}

	// Check if it's all numeric
	for _, r := range snowflake {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Additional check: Discord snowflakes shouldn't start with 0
	if snowflake[0] == '0' {
		return false
	}

	return true
}

// RuneScape Name validator
func validateRSN(fl validator.FieldLevel) bool {
	rsn := fl.Field().String()

	// Check length (RS names are 1-12 characters)
	if len(rsn) == 0 || len(rsn) > 12 {
		return false
	}

	// Check valid characters: letters, numbers, spaces, hyphens, underscores
	for _, r := range rsn {
		if !isValidRSNChar(r) {
			return false
		}
	}

	// RSN can't start or end with space
	if rsn[0] == ' ' || rsn[len(rsn)-1] == ' ' {
		return false
	}

	// Can't have consecutive spaces
	if strings.Contains(rsn, "  ") {
		return false
	}

	return true
}

func isValidRSNChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == ' ' || r == '-' || r == '_'
}

func validatePositiveTime(fl validator.FieldLevel) bool {
	return fl.Field().Int() > 0
}

// Boss name validator (checks against database)
func validateBossName(fl validator.FieldLevel) bool {
	bossName := strings.ToLower(fl.Field().String())

	return validBosses[bossName]
}

// Load valid boss names from database (call this at startup)
func loadValidBosses() {
	// This would typically be called from your main function or handler init
	// For now, we'll use a static list, but you should load from your DB
	validBosses = map[string]bool{
		"zulrah":                 true,
		"vorkath":                true,
		"cox":                    true, // Chambers of Xeric
		"tob":                    true, // Theatre of Blood
		"nightmare":              true,
		"phosani":                true,
		"nex":                    true,
		"corp":                   true, // Corporeal Beast
		"gwd":                    true, // God Wars Dungeon
		"barrows":                true,
		"dagannoth_kings":        true,
		"kalphite_queen":         true,
		"chaos_elemental":        true,
		"crazy_archaeologist":    true,
		"deranged_archaeologist": true,
		"scorpia":                true,
		"venenatis":              true,
		"vetion":                 true,
		"callisto":               true,
		// Add all your boss names here, or better yet, load from database
	}
}

// Load bosses from database (call this from your main function)
// func LoadBossesFromDB(ctx context.Context, queries interface{}) error {
// This is a placeholder - implement based on your database queries
// bosses, err := database.Queries.GetBosses(ctx)
// if err != nil {
//     return err
// }

// bossesMutex.Lock()
// validBosses = make(map[string]bool, len(bosses))
// for _, boss := range bosses {
//     validBosses[strings.ToLower(boss.Name)] = true
// }
// bossesMutex.Unlock()

// 	return nil
// }

// User-friendly validation messages
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", fe.Field(), fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
		}
		return fmt.Sprintf("%s must be at most %s", fe.Field(), fe.Param())
	case "discord_snowflake":
		return fmt.Sprintf("%s must be a valid Discord ID (17-19 digits only and can't start with a zero)", fe.Field())
	case "rsn":
		return fmt.Sprintf("%s must be a valid RuneScape name (1-12 chars, letters/numbers/spaces/hyphens/underscores only)", fe.Field())
	case "positive_time":
		return fmt.Sprintf("%s must be a positive number", fe.Field())
	case "boss_name":
		return fmt.Sprintf("%s must be a valid boss name", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "dive":
		return fmt.Sprintf("invalid item in %s array", fe.Field())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
