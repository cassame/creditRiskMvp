package internal

import (
	"context"
	"credit-risk-mvp/internal/config"

	"golang.org/x/sync/errgroup"
)

func RunStrategy(ctx context.Context, cfg config.Config, app Application, s Strategy) ([]CheckResult, error) {
	g, ctx := errgroup.WithContext(ctx)
	checks := make([]CheckResult, len(s.Checks))
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
func DecideStatus(checks []CheckResult) string {
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

func ChooseStrategy(app Application) Strategy {
	if app.Residency == "resident" && app.FirstTime {
		return Strategy{
			Name: "resident_first_time",
			Checks: []CheckFunc{
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
		return Strategy{
			Name: "resident_repeat",
			Checks: []CheckFunc{
				checkAge,
				checkPhone,
				checkPassport,
				checkAmountLimit,
			},
		}
	}
	if app.Residency == "nonresident" && app.FirstTime {
		return Strategy{
			Name: "nonresident_first_time",
			Checks: []CheckFunc{
				checkAge,
				checkPhone,
				checkPassport,
				checkAmountLimit,
				checkTerroristCF,
				checkCreditHistoryCF,
			},
		}
	}
	return Strategy{
		Name: "nonresident_repeat",
		Checks: []CheckFunc{
			checkAge,
			checkPhone,
			checkAmountLimit,
		},
	}
}
