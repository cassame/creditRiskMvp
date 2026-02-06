package internal

import (
	"context"
	"credit-risk-mvp/internal/config"
	"regexp"
	"strings"
	"time"
)

const baseAmountLimit = 100000 //base max credit amount

func approveAmount(requested int) (bool, string) {
	limit := baseAmountLimit
	if requested > limit {
		return false, "requested amount exceeds limit"
	}
	return true, ""
}

func isAdult(birthdate string) (bool, string) {
	t, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return false, "birthdate must be in format YYYY-MM-DD"
	}
	now := time.Now()
	age := now.Year() - t.Year()
	//if theres no birthdate in this year then decrease age by 1
	if now.Month() < t.Month() || (now.Month() == t.Month() && now.Day() < t.Day()) {
		age--
	}
	if age < 18 {
		return false, "client is under 18"
	}
	return true, ""
}

func isValidPhone(phone string) (bool, string) {
	re := regexp.MustCompile(`^\+?\d{10,15}$`)
	if !re.MatchString(phone) {
		return false, "invalid phone format"
	}
	return true, ""
}

func isValidPassport(passport string) (bool, string) {
	//1234 567890 or 1234567890
	re := regexp.MustCompile(`^\d{4}\s?\d{6}$`)
	if !re.MatchString(passport) {
		return false, "invalid passport format, expected '1234 567890'"
	}
	return true, ""
}

func hasPatronymic(fullName string) (bool, string) {
	parts := strings.Fields(fullName)
	if len(parts) < 3 {
		return false, "patronymic is required (expected full name with 3 parts)"
	}
	return true, ""
}

func localCheck(name string, ok bool, reason string) CheckResult {
	status := "passed"
	if !ok {
		status = "failed"
	}
	return CheckResult{Check: name, Status: status, Reason: reason}
}

func checkAge(ctx context.Context, cfg config.Config, app Application) CheckResult {
	ok, reason := isAdult(app.Birthdate)
	return localCheck("age>=18", ok, reason)
}
func checkPhone(ctx context.Context, cfg config.Config, app Application) CheckResult {
	ok, reason := isValidPhone(app.Phone)
	return localCheck("valid_phone", ok, reason)
}
func checkPassport(ctx context.Context, cfg config.Config, app Application) CheckResult {
	ok, reason := isValidPassport(app.Passport)
	return localCheck("valid_passport", ok, reason)
}
func checkPatronymic(ctx context.Context, cfg config.Config, app Application) CheckResult {
	ok, reason := hasPatronymic(app.Name)
	return localCheck("has_patronymic", ok, reason)
}
func checkAmountLimit(ctx context.Context, cfg config.Config, app Application) CheckResult {
	ok, reason := approveAmount(app.RequestedAmount)
	return localCheck("approve_amount", ok, reason)
}
