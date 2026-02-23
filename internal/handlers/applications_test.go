package handlers

import (
	"bytes"
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	savedApp domain.Application
}

func (m *mockRepo) SaveApplication(ctx context.Context, app domain.Application) error {
	m.savedApp = app
	return nil
}

// Мок очереди
type mockQueue struct{}

func (m mockQueue) SendMessage(topic string, msg []byte) error { return nil }

func TestApplicationHandler_Success(t *testing.T) {
	mockRepo := new(storage.MockRepository)
	mockQueue := new(services.MockQueue)
	mockNotify := new(notifier.MockNotifier)

	mockRepo.On("SaveApplication", mock.Anything, mock.AnythingOfType("domain.Application")).Return(nil)
	mockNotify.On("Notify", mock.Anything, mock.Anything).Return(nil)
	mockQueue.On("SendMessage", mock.Anything, mock.Anything).Return(nil)

	oldNotify := notify
	notify = mockNotify
	defer func() { notify = oldNotify }()

	cfg := config.Config{HTTPtimeout: 2 * time.Second}
	handler := &ApplicationHandler{
		Repo:  mockRepo,
		Cfg:   cfg,
		Queue: mockQueue,
	}

	payload := map[string]any{
		"name":             "Ivan Ivanov Ivanovich",
		"birthdate":        "1990-01-01",
		"phone":            "+79991234567",
		"passport":         "1234 567890",
		"residency":        "resident",
		"first_time":       true,
		"requested_amount": 5000,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertCalled(t, "SaveApplication", mock.Anything, mock.MatchedBy(func(app domain.Application) bool {
		return string(app.Name) == "Ivan Ivanov Ivanovich" && app.Status == "manual_review"
	}))
	mockNotify.AssertExpectations(t)
}
