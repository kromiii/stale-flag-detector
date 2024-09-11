package unleash

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/kromiii/stale-flag-detector/config"
)

type UnleashClient struct {
	BaseURL   string
	APIToken  string
	ProjectID string
	Config    *config.Config
}

type FeatureFlag struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	Enabled   bool      `json:"enabled"`
	Stale     bool      `json:"stale"`
}

type FeatureFlagsResponse struct {
	Features []FeatureFlag `json:"features"`
}

func NewClient(baseURL, apiToken string, projectID string, cfg *config.Config) *UnleashClient {
	return &UnleashClient{
		BaseURL:   baseURL,
		APIToken:  apiToken,
		ProjectID: projectID,
		Config:    cfg,
	}
}

func (c *UnleashClient) GetStaleFlags(excludePotentiallyStaleFlags bool) ([]string, error) {
	url := fmt.Sprintf("%s/admin/projects/%s/features", c.BaseURL, c.ProjectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp == nil {
		return nil, fmt.Errorf("API request failed: no response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var featureFlagsResp FeatureFlagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&featureFlagsResp); err != nil {
		return nil, err
	}

	return c.getStaleFlags(featureFlagsResp.Features, excludePotentiallyStaleFlags), nil
}

func (c *UnleashClient) getStaleFlags(flags []FeatureFlag, excludePotentiallyStaleFlags bool) []string {
	var staleFlags []string
	now := time.Now()

	for _, flag := range flags {
		if c.isFlagStale(flag, now, excludePotentiallyStaleFlags) {
			staleFlags = append(staleFlags, flag.Name)
		}
	}

	return staleFlags
}

func (c *UnleashClient) isFlagStale(flag FeatureFlag, now time.Time, excludePotentiallyStaleFlags bool) bool {
	if excludePotentiallyStaleFlags || flag.Stale {
		return flag.Stale
	}
	lifetime := c.getExpectedLifetime(flag.Type)
	return flag.Stale || now.Sub(flag.CreatedAt) > lifetime
}

func (c *UnleashClient) getExpectedLifetime(flagType string) time.Duration {
	lifetimeMap := map[string]int{
		"release":     c.Config.ReleaseFlagLifetime,
		"experiment":  c.Config.ExperimentFlagLifetime,
		"operational": c.Config.OperationalFlagLifetime,
		"permission":  c.Config.PermissionFlagLifetime,
	}

	if lifetime, ok := lifetimeMap[flagType]; ok {
		if lifetime == -1 {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(lifetime) * 24 * time.Hour
	}

	if flagType == "kill-switch" {
		return time.Duration(math.MaxInt64)
	}

	return 30 * 24 * time.Hour // デフォルトは30日
}
