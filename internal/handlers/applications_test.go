package handlers

import (
	"bytes"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
	"credit-risk-mvp/internal/logger"
	"credit-risk-mvp/internal/repository"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.InitLogger("local")

	os.Exit(m.Run())
}

func TestApplicationHandler_Success(t *testing.T) {
	mockRepo := new(storage.MockRepository)
	mockQueue := new(services.MockQueue)
	mockNotify := new(notifier.MockNotifier)

	mockRepo.On("SaveApplication", mock.Anything, mock.AnythingOfType("domain.Application")).Return(nil)
	mockNotify.On("Notify", mock.Anything, mock.Anything).Return(nil)
	mockQueue.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	cfg := config.Config{HTTPtimeout: 2 * time.Second}
	handler := &ApplicationHandler{
		Repo:           mockRepo,
		Cfg:            cfg,
		Queue:          mockQueue,
		Notifier:       mockNotify,
		TerroristStore: repository.NewLocalTerroristStore(),
	}

	payload := map[string]any{
		"name":             "Ivan Ivanov Ivanovich",
		"birthdate":        "1990-01-01",
		"phone":            "+79991234567",
		"passport":         "1234 567890",
		"residency":        "resident",
		"first_time":       true,
		"requested_amount": 50000,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertCalled(t, "SaveApplication", mock.Anything, mock.MatchedBy(func(app domain.Application) bool {
		return string(app.Name) == "Ivan Ivanov Ivanovich"
	}))
	mockNotify.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}
