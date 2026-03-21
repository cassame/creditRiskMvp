package domain

import (
	"context"
	"credit-risk-mvp/internal/config"
	"time"
)

type Application struct {
	//Lifecycle
	ID           string        `json:"application_id"`
	Status       string        `json:"status"`
	StrategyName string        `json:"strategy"`
	Checks       []CheckResult `json:"checks"`
	CreatedAt    time.Time     `json:"created_at"`
	//Client data
	Payload         map[string]any `json:"payload"`
	Name            FullName       `json:"-"`
	Birthdate       Birthdate      `json:"-"`
	Phone           Phone          `json:"-"`
	Passport        Passport       `json:"-"`
	RequestedAmount Amount         `json:"-"`
	Residency       string         `json:"-"`
	FirstTime       bool           `json:"-"`
}

type CheckFunc func(ctx context.Context, cfg config.Config, app Application) CheckResult

type CheckResult struct {
	Check  string `json:"check"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type Strategy struct {
	Name   string
	Checks []CheckFunc
}
