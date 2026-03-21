package internal

import (
	"context"
	"credit-risk-mvp/internal/config"
	"credit-risk-mvp/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group
var lastTerroristCacheUpdate time.Time

func checkCreditHistoryCF(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	return checkCreditHistory(ctx, cfg, string(app.Passport))
}

// checkCreditHistory checking for credit history score (passed/failed/error)
func checkCreditHistory(ctx context.Context, cfg config.Config, passport string) domain.CheckResult {

	client := &http.Client{Timeout: cfg.HTTPtimeout}
	u := cfg.CreditHistoryURL + "?passport=" + url.QueryEscape(passport)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return domain.CheckResult{
			Check:  "credit_history",
			Status: "error",
			Reason: "failed to create request: " + err.Error(),
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return domain.CheckResult{
			Check:  "credit_history",
			Status: "error",
			Reason: "external service unavailable: " + err.Error(),
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return domain.CheckResult{
			Check:  "credit_history",
			Status: "error",
			Reason: "external service bad status: " + resp.Status,
		}
	}
	var data struct {
		IsGood bool `json:"is_good"`
		Score  int  `json:"score"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return domain.CheckResult{
			Check:  "credit_history",
			Status: "error",
			Reason: "external service bad json: " + err.Error(),
		}
	}
	if !data.IsGood {
		return domain.CheckResult{
			Check:  "credit_history",
			Status: "failed",
			Reason: "bad credit history",
		}
	}
	return domain.CheckResult{
		Check:  "credit_history",
		Status: "passed",
	}
}

// NewTerroristChecker adapter to return CheckResult from checkTerrorist
func NewTerroristChecker(store domain.TerroristStore) func(context.Context, config.Config, domain.Application) domain.CheckResult {
	return func(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
		return checkTerrorist(ctx, cfg, store, string(app.Passport))
	}
}

// checkTerrorist checking for terrorist (passed/failed/error)
func checkTerrorist(ctx context.Context, cfg config.Config, store domain.TerroristStore, passport string) domain.CheckResult {
	needRefresh := lastTerroristCacheUpdate.IsZero() || time.Since(lastTerroristCacheUpdate) > time.Hour*24

	if needRefresh {
		if err := refreshTerroristCache(ctx, cfg, store); err != nil {
			return domain.CheckResult{
				Check:  "terrorist",
				Status: "error",
				Reason: "cannot refresh terrorist list: " + err.Error(),
			}
		}
	}
	found, err := store.IsTerrorist(ctx, passport)
	if err != nil {
		return domain.CheckResult{Check: "terrorist", Status: "error", Reason: err.Error()}
	}
	if found {
		return domain.CheckResult{
			Check:  "terrorist",
			Status: "failed",
			Reason: "client is in terrorist/extremist list",
		}
	}
	return domain.CheckResult{Check: "terrorist", Status: "passed"}
}

func actualRefresh(ctx context.Context, cfg config.Config, store domain.TerroristStore) error {
	client := &http.Client{Timeout: cfg.HTTPtimeout}

	req, err := http.NewRequestWithContext(ctx, "GET", cfg.TerroristURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	var data struct {
		Passports []string `json:"passports"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	return store.UpdateList(ctx, data.Passports)
}

// refreshTerroristCache updating info from TerroristURL and updating time
func refreshTerroristCache(ctx context.Context, cfg config.Config, store domain.TerroristStore) error {
	_, err, _ := sf.Do("refresh", func() (interface{}, error) {
		err := actualRefresh(ctx, cfg, store)
		if err != nil {
			return nil, err
		}
		lastTerroristCacheUpdate = time.Now()
		return nil, nil
	})
	return err
}

func checkBankruptcy(ctx context.Context, cfg config.Config, passport string) domain.CheckResult {
	client := &http.Client{Timeout: 1 * cfg.HTTPtimeout}

	u, err := url.Parse(cfg.BankruptcyURL)
	if err != nil {
		return domain.CheckResult{
			Check:  "bankruptcy",
			Status: "error",
			Reason: "invalid URL in config",
		}
	}
	q := u.Query()
	q.Set("passport", passport)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return domain.CheckResult{Check: "bankruptcy", Status: "error", Reason: "failed to create request"}
	}
	resp, err := client.Do(req)
	if err != nil {
		return domain.CheckResult{Check: "bankruptcy", Status: "error", Reason: "..."}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return domain.CheckResult{Check: "bankruptcy", Status: "error", Reason: "bad status"}
	}
	var data struct {
		IsBankrupt bool `json:"is_bankrupt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return domain.CheckResult{Check: "bankruptcy", Status: "error", Reason: "bad json"}
	}
	if data.IsBankrupt {
		return domain.CheckResult{Check: "bankruptcy", Status: "failed", Reason: "client is bankrupt"}
	}
	return domain.CheckResult{Check: "bankruptcy", Status: "passed"}
}
