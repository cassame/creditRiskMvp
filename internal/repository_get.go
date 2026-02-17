package internal

import (
	"database/sql"
	"encoding/json"
	"time"
)

type ApplicationView struct {
	ApplicationID string         `json:"application_id"`
	CreatedAt     time.Time      `json:"created_at"`
	Strategy      string         `json:"strategy"`
	Status        string         `json:"status"`
	Payload       map[string]any `json:"payload"`
	Checks        []CheckResult  `json:"checks"`
}

func GetApplication(db *sql.DB, id string) (ApplicationView, error) {
	var (
		strategy string
		status   string
		created  time.Time
		payloadB []byte
	)
	err := db.QueryRow(`
		SELECT strategy, status, created_at, payload
		FROM applications
		WHERE id = $1
		`, id).Scan(&strategy, &status, &created, &payloadB)
	if err != nil {
		return ApplicationView{}, err
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadB, &payload); err != nil {
		return ApplicationView{}, err
	}
	rows, err := db.Query(`
		SELECT check_name, status, reason
		FROM check_results
		WHERE application_id = $1
		order by id asc`, id)
	checks := make([]CheckResult, 0)
	for rows.Next() {
		var c CheckResult
		if err := rows.Scan(&c.Check, &c.Status, &c.Reason); err != nil {
			return ApplicationView{}, err
		}
		checks = append(checks, c)
		if err := rows.Err(); err != nil {
			return ApplicationView{}, err
		}
	}
	return ApplicationView{
		ApplicationID: id,
		CreatedAt:     created,
		Strategy:      strategy,
		Status:        status,
		Payload:       payload,
		Checks:        checks,
	}, nil
}
