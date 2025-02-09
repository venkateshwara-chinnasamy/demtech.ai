package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"mock-ses-api/internal/handlers"
	"mock-ses-api/internal/models"
	"mock-ses-api/pkg/stats"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := handlers.NewSESHandler(stats.NewEmailStats())
	router.GET("/health", handler.HealthCheck)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// Serve request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name           string
		request        models.SendEmailRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid request",
			request: models.SendEmailRequest{
				Destination: models.Destination{
					ToAddresses: []string{"test@example.com"},
				},
				Message: models.Message{
					Subject: models.Subject{
						Data: "Test Subject",
					},
					Body: models.Body{
						Text: models.Html{
							Data: "Test Body",
						},
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "invalid request",
			request: models.SendEmailRequest{
				Destination: models.Destination{},
				Message: models.Message{
					Subject: models.Subject{
						Data: "Test Subject",
					},
					Body: models.Body{
						Text: models.Html{
							Data: "Test Body",
						},
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()

			handler := handlers.NewSESHandler(stats.NewEmailStats())
			router.POST("/outbound-emails", handler.SendEmail)

			// Create request
			body, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/outbound-emails", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Serve request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectedError {
				var response models.SendEmailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.MessageId)
				assert.NotEmpty(t, response.RequestId)
			}

		})
	}
}
