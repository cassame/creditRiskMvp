package storage

import (
	"credit-risk-mvp/internal"
	"database/sql"
	"encoding/json"
)

func SaveApplication(db *sql.DB, appID string, strategyName string, status string, payload map[string]any, checks []internal.CheckResult) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	//APPLICATIONS INSERT
	_, err = tx.Exec(`
		INSERT INTO APPLICATIONS (id, strategy, status, payload)
		VALUES ($1, $2, $3, $4)
		`, appID, strategyName, status, payloadBytes)
	if err != nil {
		return err
	}
	//CHECK_RESULTS INSERT
	for _, c := range checks {
		_, err = tx.Exec(`
			INSERT INTO check_results (application_id, check_name, status, reason)
			VALUES ($1, $2, $3, $4)`, appID, c.Check, c.Status, c.Reason)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
