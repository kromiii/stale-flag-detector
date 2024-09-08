package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	UnleashAPIEndpoint string
	UnleashAPIToken    string
	ProjectID          string
	OpenAIAPIKey       string
	GitHubToken        string
	GitHubOwner        string
	GitHubRepo         string
	GitHubBaseURL			string
	ReleaseFlagLifetime int
	ExperimentFlagLifetime int
	OperationalFlagLifetime int
	PermisionFlagLifetime int
}

func Load() (*Config, error) {
	endpoint := os.Getenv("UNLEASH_API_ENDPOINT")
	token := os.Getenv("UNLEASH_API_TOKEN")
	projectID := os.Getenv("UNLEASH_PROJECT_ID")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubBaseURL := os.Getenv("GITHUB_BASE_URL")
	releaseFlagLifetime := os.Getenv("RELEASE_FLAG_LIFETIME")
	experimentFlagLifetime := os.Getenv("EXPERIMENT_FLAG_LIFETIME")
	operationalFlagLifetime := os.Getenv("OPERATIONAL_FLAG_LIFETIME")
	permissionFlagLifetime := os.Getenv("PERMISSION_FLAG_LIFETIME")

	if endpoint == "" || token == "" || openaiKey == "" || projectID == "" || githubToken == "" || githubOwner == "" || githubRepo == "" {
		return nil, errors.New("missing required environment variables")
	}

	releaseFlagLifetimeInt, err := parseLifetime(releaseFlagLifetime)
	if err != nil {
			return nil, fmt.Errorf("invalid release flag lifetime: %v", err)
	}

	experimentFlagLifetimeInt, err := parseLifetime(experimentFlagLifetime)
	if err != nil {
			return nil, fmt.Errorf("invalid experiment flag lifetime: %v", err)
	}

	operationalFlagLifetimeInt, err := parseLifetime(operationalFlagLifetime)
	if err != nil {
			return nil, fmt.Errorf("invalid operational flag lifetime: %v", err)
	}

	permissionFlagLifetimeInt, err := parseLifetime(permissionFlagLifetime)
	if err != nil {
			return nil, fmt.Errorf("invalid permission flag lifetime: %v", err)
	}

	return &Config{
		UnleashAPIEndpoint: endpoint,
		UnleashAPIToken:    token,
		ProjectID:          projectID,
		OpenAIAPIKey:       openaiKey,
		GitHubToken:        githubToken,
		GitHubOwner:        githubOwner,
		GitHubRepo:         githubRepo,
		GitHubBaseURL:      githubBaseURL,
		ReleaseFlagLifetime: releaseFlagLifetimeInt,
		ExperimentFlagLifetime: experimentFlagLifetimeInt,
		OperationalFlagLifetime: operationalFlagLifetimeInt,
		PermisionFlagLifetime: permissionFlagLifetimeInt,
	}, nil
}

func parseLifetime(lifetime string) (int, error) {
	if lifetime == "permanent" {
			return -1, nil
	}
	if lifetime == "" {
		return 30, nil // default to 30 days
	}
	return strconv.Atoi(lifetime)
}
