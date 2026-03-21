package internal

import (
	"context"
	"credit-risk-mvp/internal/config"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckBankruptcy_Mock(t *testing.T) {
	tests := []struct {
		name           string
		mockResponce   string
		mockStatus     int
		expectedStatus string
		delay          time.Duration
	}{
		{
			name:           "Client is bankrupt",
			mockResponce:   `{"is_bankrupt": true}`,
			mockStatus:     http.StatusOK,
			expectedStatus: "failed",
			delay:          0 * time.Second,
		},
		{
			name:           "Client is not bankrupt",
			mockResponce:   `{"is_bankrupt": false}`,
			mockStatus:     http.StatusOK,
			expectedStatus: "passed",
			delay:          0 * time.Second,
		},
		{
			name:           "External service error",
			mockResponce:   `Internal Server Error`,
			mockStatus:     http.StatusInternalServerError,
			expectedStatus: "error",
			delay:          0 * time.Second,
		},
		{
			name:           "Timeout case",
			mockResponce:   `{"is_bankrupt": false}`,
			mockStatus:     http.StatusOK,
			expectedStatus: "error",
			delay:          2 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.delay > 0 {
					time.Sleep(tt.delay)
				}
				w.WriteHeader(tt.mockStatus)
				_, err := fmt.Fprint(w, tt.mockResponce)
				if err != nil {
					t.Fatalf("failed: %s", err.Error())
				}
			}))
			defer server.Close()
			cfg := config.Config{
				BankruptcyURL: server.URL,
				HTTPtimeout:   1 * time.Second,
			}
			res := checkBankruptcy(context.Background(), cfg, "1234567890")
			if res.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s (reason: %s)",
					tt.expectedStatus, res.Status, res.Reason)
			}
		})
	}
}
