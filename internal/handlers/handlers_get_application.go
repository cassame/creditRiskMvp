package handlers

import (
	"credit-risk-mvp/internal"
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

func MakeGetApplicationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		//path: /applications/{id}
		id := strings.TrimPrefix(r.URL.Path, "/applications/")
		if id == "" || id == r.URL.Path {
			writeError(w, http.StatusBadRequest, "application id is required")
			return
		}

		appView, err := internal.GetApplication(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot load application")
			return
		}
		writeJSON(w, http.StatusOK, appView)
	}
}
