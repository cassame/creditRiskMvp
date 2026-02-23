package internal

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"

	"golang.org/x/sync/errgroup"
)

func RunStrategy(ctx context.Context, cfg config.Config, app domain.Application, s domain.Strategy) ([]domain.CheckResult, error) {
	g, ctx := errgroup.WithContext(ctx)
	checks := make([]domain.CheckResult, len(s.Checks))
	for i, fn := range s.Checks {
		i, fn := i, fn
		g.Go(func() error {
			checks[i] = fn(ctx, cfg, app)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return checks, nil
}
func DecideStatus(checks []domain.CheckResult) string {
	hasCriticalError := false

	for _, c := range checks {
		if c.Status == "failed" {
			return "rejected"
		}
		if c.Status == "error" && (c.Check == "bankruptcy" || c.Check == "terrorist") {
			hasCriticalError = true
		}
	}
	if hasCriticalError {
		return "manual_review"
	}
	return "approved"
}

func ChooseStrategy(app domain.Application) domain.Strategy {
	if app.Residency == "resident" && app.FirstTime {
		return domain.Strategy{
			Name: "resident_first_time",
			Checks: []domain.CheckFunc{
				checkAge,
				checkPhone,
				checkPassport,
				checkPatronymic,
				checkAmountLimit,
				checkTerroristCF,
				checkCreditHistoryCF,
			},
		}
	}
	if app.Residency == "resident" && !app.FirstTime {
		return domain.Strategy{
			Name: "resident_repeat",
			Checks: []domain.CheckFunc{
				checkAge,
				checkPhone,
				checkPassport,
				checkAmountLimit,
			},
		}
	}
	if app.Residency == "nonresident" && app.FirstTime {
		return domain.Strategy{
			Name: "nonresident_first_time",
			Checks: []domain.CheckFunc{
				checkAge,
				checkPhone,
				checkPassport,
				checkAmountLimit,
				checkTerroristCF,
				checkCreditHistoryCF,
			},
		}
	}
	return domain.Strategy{
		Name: "nonresident_repeat",
		Checks: []domain.CheckFunc{
			checkAge,
			checkPhone,
			checkAmountLimit,
		},
	}
}
