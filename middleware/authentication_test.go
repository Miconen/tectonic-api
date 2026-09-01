package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"tectonic-api/config"
	"tectonic-api/logging"
	"tectonic-api/models"
)

const (
	testAPIKey     = "full-access-key"
	testAPIReadKey = "read-only-key"
)

func TestMain(m *testing.M) {
	logging.Init(&config.Config{LogLevel: "error"})
	os.Exit(m.Run())
}

func TestAuthentication(t *testing.T) {
	cfg := &config.Config{
		APIKey:     testAPIKey,
		APIReadKey: testAPIReadKey,
	}

	tests := []struct {
		name           string
		method         string
		path           string
		token          string
		expectedStatus int
		expectedError  *models.APIV1ErrorCode
		handlerCalled  bool
	}{
		{
			name:           "full key can read",
			method:         http.MethodGet,
			path:           "/api/v1/guilds/123",
			token:          testAPIKey,
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "full key can post",
			method:         http.MethodPost,
			path:           "/api/v1/guilds",
			token:          testAPIKey,
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "full key can put",
			method:         http.MethodPut,
			path:           "/api/v1/guilds/123",
			token:          testAPIKey,
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "full key can delete",
			method:         http.MethodDelete,
			path:           "/api/v1/guilds/123",
			token:          testAPIKey,
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "read key can get",
			method:         http.MethodGet,
			path:           "/api/v1/guilds/123",
			token:          testAPIReadKey,
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "read key can use head",
			method:         http.MethodHead,
			path:           "/api/v1/guilds/123",
			token:          testAPIReadKey,
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "read key cannot post",
			method:         http.MethodPost,
			path:           "/api/v1/guilds",
			token:          testAPIReadKey,
			expectedStatus: http.StatusForbidden,
			expectedError:  errorCode(models.ERROR_INSUFFICIENT_SCOPE),
		},
		{
			name:           "read key cannot put",
			method:         http.MethodPut,
			path:           "/api/v1/guilds/123",
			token:          testAPIReadKey,
			expectedStatus: http.StatusForbidden,
			expectedError:  errorCode(models.ERROR_INSUFFICIENT_SCOPE),
		},
		{
			name:           "read key cannot delete",
			method:         http.MethodDelete,
			path:           "/api/v1/guilds/123",
			token:          testAPIReadKey,
			expectedStatus: http.StatusForbidden,
			expectedError:  errorCode(models.ERROR_INSUFFICIENT_SCOPE),
		},
		{
			name:           "invalid key is unauthorized",
			method:         http.MethodGet,
			path:           "/api/v1/guilds/123",
			token:          "invalid-key",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  errorCode(models.ERROR_INVALID_TOKEN),
		},
		{
			name:           "missing key is unauthorized",
			method:         http.MethodGet,
			path:           "/api/v1/guilds/123",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  errorCode(models.ERROR_INVALID_TOKEN),
		},
		{
			name:           "docs remain public",
			method:         http.MethodGet,
			path:           "/docs",
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
		{
			name:           "openapi spec remains public",
			method:         http.MethodGet,
			path:           "/openapi.json",
			expectedStatus: http.StatusNoContent,
			handlerCalled:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			called := false

			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			})

			request := httptest.NewRequest(test.method, test.path, nil)
			if test.token != "" {
				request.Header.Set("Authorization", test.token)
			}

			response := httptest.NewRecorder()
			Authentication(cfg)(next).ServeHTTP(response, request)

			if response.Code != test.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					test.expectedStatus,
					response.Code,
				)
			}

			if called != test.handlerCalled {
				t.Fatalf(
					"expected downstream handler called=%t, got %t",
					test.handlerCalled,
					called,
				)
			}

			if test.expectedError == nil {
				return
			}

			var apiError models.TectonicError
			if err := json.NewDecoder(response.Body).Decode(&apiError); err != nil {
				t.Fatalf("decode error response: %v", err)
			}

			if apiError.ErrCode != test.expectedError.Code() {
				t.Errorf(
					"expected error code %d, got %d",
					test.expectedError.Code(),
					apiError.ErrCode,
				)
			}

			if apiError.Msg != test.expectedError.Message() {
				t.Errorf(
					"expected error message %q, got %q",
					test.expectedError.Message(),
					apiError.Msg,
				)
			}

			if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf(
					"expected Content-Type application/json, got %q",
					contentType,
				)
			}
		})
	}
}

func errorCode(code models.APIV1ErrorCode) *models.APIV1ErrorCode {
	return &code
}
