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

var tcache terroristCache

// var tmu sync.Mutex
var sf singleflight.Group

type terroristCache struct {
	lastUpdated time.Time
	set         map[string]struct{}
}

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

// checkTerroristCF adapter to return CheckResult from checkTerrorist
func checkTerroristCF(ctx context.Context, cfg config.Config, app domain.Application) domain.CheckResult {
	return checkTerrorist(ctx, cfg, string(app.Passport))
}

// checkTerrorist checking for terrorist (passed/failed/error)
func checkTerrorist(ctx context.Context, cfg config.Config, passport string) domain.CheckResult {
	needRefresh := tcache.set == nil || time.Since(tcache.lastUpdated) > time.Hour*24
	if needRefresh {
		if err := refreshTerroristCache(ctx, cfg); err != nil {
			return domain.CheckResult{
				Check:  "terrorist",
				Status: "error",
				Reason: "cannot refresh terrorist list: " + err.Error(),
			}
		}
	}
	_, found := tcache.set[passport]
	if found {
		return domain.CheckResult{
			Check:  "terrorist",
			Status: "failed",
			Reason: "client is in terrorist/extremist list",
		}
	}
	return domain.CheckResult{Check: "terrorist", Status: "passed"}
}

func actualRefresh(ctx context.Context, cfg config.Config) error {
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
	newSet := make(map[string]struct{})
	for _, passport := range data.Passports {
		newSet[passport] = struct{}{}
	}
	tcache.set = newSet
	tcache.lastUpdated = time.Now()
	return nil

}

// refreshTerroristCache updating info from TerroristURL and updating time
func refreshTerroristCache(ctx context.Context, cfg config.Config) error {
	_, err, _ := sf.Do("refresh", func() (interface{}, error) {
		return nil, actualRefresh(ctx, cfg)
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
