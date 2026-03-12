package handlers

import (
	"credit-risk-mvp/internal"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ApplicationHandler struct {
	Repo     domain.Repository
	Cfg      config.Config
	Queue    services.MessageQueue
	Notifier notifier.Notifier
}

func NewApplicationsHandler(cfg config.Config, repo domain.Repository, queue services.MessageQueue, n notifier.Notifier) *ApplicationHandler {
	return &ApplicationHandler{
		Repo:     repo,
		Cfg:      cfg,
		Queue:    queue,
		Notifier: n,
	}
}

func (h *ApplicationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	app, err := internal.ParseApplication(payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	strategy := internal.ChooseStrategy(app)
	app.StrategyName = strategy.Name
	app.Checks, err = internal.RunStrategy(ctx, h.Cfg, app, strategy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to run strategy")
		return
	}
	app.Status = internal.DecideStatus(app.Checks)
	app.ID = uuid.NewString()

	if err := h.Repo.SaveApplication(ctx, app); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot save application")
		return
	}

	message := fmt.Sprintf("New application: %s with status: %s", app.ID, app.Status)
	err = h.Queue.SendMessage("application_topic", app.ID, []byte(message))

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send message to Kafka")
		return
	}
	err = h.Notifier.Notify(app, app.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to notify application")
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"application_id": app.ID,
		"status":         app.Status,
		"strategy":       strategy.Name,
		"checks":         app.Checks,
	})

}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
