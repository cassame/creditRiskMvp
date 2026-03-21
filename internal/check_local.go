package internal

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
)

const baseAmountLimit = 100000 //base max credit amount

func localCheck(name string, ok bool, reason string) domain.CheckResult {
	status := "passed"
	if !ok {
		return domain.CheckResult{Check: name, Status: "failed", Reason: reason}
	}
	return domain.CheckResult{Check: name, Status: status, Reason: ""}
}

func checkAge(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	age := app.Birthdate.Age()
	ok := age >= 18
	reason := ""
	if !ok {
		reason = "client is under 18"
	}
	return localCheck("age>=18", ok, reason)
}
func checkPhone(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	return localCheck("valid_phone", true, "")
}
func checkPassport(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	return localCheck("valid_passport", true, "")
}
func checkPatronymic(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	ok := app.Name.HasPatronymic()
	return localCheck("has_patronymic", ok, "patronymic is required")
}
func checkAmountLimit(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	ok := int(app.RequestedAmount) <= baseAmountLimit
	reason := ""
	if !ok {
		reason = "requested amount exceeds limit"
	}
	return localCheck("approve_amount", ok, reason)
}
