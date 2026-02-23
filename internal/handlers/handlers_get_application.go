package handlers

import (
	"credit-risk-mvp/internal/domain"
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

func MakeGetApplicationHandler(repo domain.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/applications/")
		if id == "" || id == r.URL.Path {
			writeError(w, http.StatusBadRequest, "application id is required")
			return
		}

		app, err := repo.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "application not found")
			} else {
				writeError(w, http.StatusInternalServerError, "db error")
			}
			return
		}
		writeJSON(w, http.StatusOK, app)
	}
}
