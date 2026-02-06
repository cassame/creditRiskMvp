package handlers

import (
	"context"
	"credit-risk-mvp/internal"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/notifier"
	"credit-risk-mvp/services"
	"credit-risk-mvp/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var notify notifier.Notifier = notifier.LogNotifier{}

func MakeApplicationsHandler(cfg config.Config, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		handleApplications(w, r, cfg, db, ctx)
	}
}

func handleApplications(w http.ResponseWriter, r *http.Request, cfg config.Config, db *sql.DB, ctx context.Context) {
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
	checks, err := internal.RunStrategy(ctx, cfg, app, strategy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to run strategy")
		return
	}
	status := internal.DecideStatus(checks)
	appID := uuid.NewString()
	if err := storage.SaveApplication(db, appID, strategy.Name, status, app.Payload, checks); err != nil {
		writeError(w, http.StatusInternalServerError, "cannot save application")
		return
	}
	message := fmt.Sprintf("New application: %s with status: %s", appID, status)
	err = services.SendMessage("application_topic", []byte(message))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send message to Kafka")
		return
	}
	err = notify.Notify(app, status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to notify application")
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"application_id": appID,
		"status":         status,
		"strategy":       strategy.Name,
		"checks":         checks,
	})
}

// writeError sending an error in format 400 {"error": "..." }
func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// writeJSON sending to user the http response CODE
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
