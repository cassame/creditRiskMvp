package internal

import (
	"context"
	"credit-risk-mvp/internal/config"
)

type CheckFunc func(ctx context.Context, cfg config.Config, app Application) CheckResult

type Application struct {
	Payload         map[string]any
	Name            string
	Birthdate       string
	Phone           string
	Passport        string
	RequestedAmount int
	Residency       string
	FirstTime       bool
}

type CheckResult struct {
	Check  string `json:"check"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type Strategy struct {
	Name   string
	Checks []CheckFunc
}
