package handlers

import (
	"credit-risk-mvp/internal"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
	"credit-risk-mvp/internal/logger"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ApplicationHandler struct {
	Repo           domain.Repository
	Cfg            config.Config
	Queue          services.MessageQueue
	Notifier       notifier.Notifier
	TerroristStore domain.TerroristStore
}

func NewApplicationsHandler(cfg config.Config, repo domain.Repository,
	queue services.MessageQueue, n notifier.Notifier, terrStore domain.TerroristStore) *ApplicationHandler {
	return &ApplicationHandler{
		Repo:           repo,
		Cfg:            cfg,
		Queue:          queue,
		Notifier:       n,
		TerroristStore: terrStore,
	}
}

func (h *ApplicationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Lg.Info("incoming request",
		"method", r.Method, "path", r.URL.Path,
	)

	ctx := r.Context()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logAndWriteError(w, "invalid json", http.StatusBadRequest, err)
		return
	}
	app, err := internal.ParseApplication(payload)
	if err != nil {
		h.logAndWriteError(w, "validation failed", http.StatusBadRequest, err)
		return
	}
	strategy := internal.ChooseStrategy(app, h.TerroristStore)
	app.StrategyName = strategy.Name
	app.Checks, err = internal.RunStrategy(ctx, h.Cfg, app, strategy)
	if err != nil {
		h.logAndWriteError(w, "failed to run strategy", http.StatusInternalServerError, err)
		return
	}
	app.Status = internal.DecideStatus(app.Checks)
	app.ID = uuid.NewString()

	if err := h.Repo.SaveApplication(ctx, app); err != nil {
		h.logAndWriteError(w, "failed to save application",
			http.StatusInternalServerError, err, "app_id", app.ID)
		return
	}

	message := fmt.Sprintf("New application: %s with status: %s", app.ID, app.Status)
	err = h.Queue.SendMessage("application_topic", app.ID, []byte(message))

	if err != nil {
		h.logAndWriteError(w, "failed to send message to Kafka", http.StatusInternalServerError,
			err, "app_id", app.ID)
		return
	}
	err = h.Notifier.Notify(app, app.Status)
	if err != nil {
		h.logAndWriteError(w, "failed to notify application", http.StatusInternalServerError,
			err, "app_id", app.ID)
	}

	logger.Lg.Info("application processed successfully",
		"app_id", app.ID,
		"status", app.Status,
	)

	writeJSON(w, http.StatusOK, map[string]any{
		"application_id": app.ID,
		"status":         app.Status,
		"strategy":       strategy.Name,
		"checks":         app.Checks,
	})

}

func (h *ApplicationHandler) logAndWriteError(w http.ResponseWriter, msg string, code int, err error, args ...any) {
	fullArgs := append([]any{"error", err}, args...)

	if code >= 500 {
		logger.Lg.Error(msg, fullArgs...)
	} else {
		logger.Lg.Warn(msg, fullArgs...)
	}

	writeError(w, code, msg)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
