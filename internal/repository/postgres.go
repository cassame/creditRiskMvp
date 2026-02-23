package repository

import (
	"context"
	"credit-risk-mvp/internal/domain"
	"database/sql"
	"encoding/json"
)

type SqlRepository struct {
	db *sql.DB
}

func NewSqlRepository(db *sql.DB) *SqlRepository {
	return &SqlRepository{db: db}
}

func (r *SqlRepository) GetByID(ctx context.Context, id string) (domain.Application, error) {

	var app domain.Application
	var payloadB []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT strategy, status, created_at, payload
		FROM applications
		WHERE id = $1
		`, id).
		Scan(&app.StrategyName, &app.Status, &app.CreatedAt, &payloadB)

	if err != nil {
		return domain.Application{}, err
	}

	app.ID = id
	if err := json.Unmarshal(payloadB, &app.Payload); err != nil {
		return domain.Application{}, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT check_name, status, reason
		FROM check_results
		WHERE application_id = $1
		ORDER BY id ASC`, id)
	if err != nil {
		return app, nil
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.CheckResult
		if err := rows.Scan(&c.Check, &c.Status, &c.Reason); err == nil {
			app.Checks = append(app.Checks, c)
		}
	}
	return app, nil
}

func (r *SqlRepository) SaveApplication(ctx context.Context, app domain.Application) error {
	payloadBytes, err := json.Marshal(app.Payload)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	//APPLICATIONS INSERT
	_, err = tx.ExecContext(ctx, `
		INSERT INTO APPLICATIONS (id, strategy, status, payload)
		VALUES ($1, $2, $3, $4)`,
		app.ID, app.StrategyName, app.Status, payloadBytes)
	if err != nil {
		return err
	}
	//CHECK_RESULTS INSERT
	for _, c := range app.Checks {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO check_results (application_id, check_name, status, reason)
			VALUES ($1, $2, $3, $4)`, app.ID, c.Check, c.Status, c.Reason)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
