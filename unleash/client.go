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

func (c *UnleashClient) GetStaleFlags(onlyStaleFlags bool) ([]string, error) {
	url := fmt.Sprintf("%s/admin/projects/%s/features", c.BaseURL, c.ProjectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var featureFlagsResp FeatureFlagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&featureFlagsResp); err != nil {
		return nil, err
	}

	return c.getStaleFlags(featureFlagsResp.Features, onlyStaleFlags), nil
}

func (c *UnleashClient) getStaleFlags(flags []FeatureFlag, onlyStaleFlags bool) []string {
	var staleFlags []string
	now := time.Now()

	for _, flag := range flags {
		lifetime := c.getExpectedLifetime(flag.Type)
		if onlyStaleFlags {
			if flag.Stale {
				staleFlags = append(staleFlags, flag.Name)
			}
		} else {
			isStale := flag.Stale || now.Sub(flag.CreatedAt) > lifetime
			if isStale {
				staleFlags = append(staleFlags, flag.Name)
			}
		}
	}

	return staleFlags
}

func (c *UnleashClient) getExpectedLifetime(flagType string) time.Duration {
	switch flagType {
	case "release":
		if c.Config.ReleaseFlagLifetime == -1 {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(c.Config.ReleaseFlagLifetime) * 24 * time.Hour
	case "experiment":
		if c.Config.ExperimentFlagLifetime == -1 {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(c.Config.ExperimentFlagLifetime) * 24 * time.Hour
	case "operational":
		if c.Config.OperationalFlagLifetime == -1 {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(c.Config.OperationalFlagLifetime) * 24 * time.Hour
	case "permission":
		if c.Config.PermissionFlagLifetime == -1 {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(c.Config.PermissionFlagLifetime) * 24 * time.Hour
	case "kill-switch":
		return time.Duration(math.MaxInt64) // kill-switchは期限切れにならない
	default:
		return 30 * 24 * time.Hour // デフォルトは30日
	}
}
