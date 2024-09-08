package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	UnleashAPIEndpoint      string
	UnleashAPIToken         string
	ProjectID               string
	ReleaseFlagLifetime     int
	ExperimentFlagLifetime  int
	OperationalFlagLifetime int
	PermissionFlagLifetime  int
}

func Load() (*Config, error) {
	config := &Config{
		UnleashAPIEndpoint: os.Getenv("UNLEASH_API_ENDPOINT"),
		UnleashAPIToken:    os.Getenv("UNLEASH_API_TOKEN"),
		ProjectID:          getEnvWithDefault("UNLEASH_PROJECT_ID", "default"),
	}

	if config.UnleashAPIEndpoint == "" || config.UnleashAPIToken == "" {
		return nil, errors.New("missing required environment variables")
	}

	lifetimes := map[string]*int{
		"RELEASE_FLAG_LIFETIME":     &config.ReleaseFlagLifetime,
		"EXPERIMENT_FLAG_LIFETIME":  &config.ExperimentFlagLifetime,
		"OPERATIONAL_FLAG_LIFETIME": &config.OperationalFlagLifetime,
		"PERMISSION_FLAG_LIFETIME":  &config.PermissionFlagLifetime,
	}

	defaultValues := map[string]string{
		"RELEASE_FLAG_LIFETIME":     "40",
		"EXPERIMENT_FLAG_LIFETIME":  "40",
		"OPERATIONAL_FLAG_LIFETIME": "7",
		"PERMISSION_FLAG_LIFETIME":  "permanent",
	}

	for envVar, lifetimePtr := range lifetimes {
		value := getEnvWithDefault(envVar, defaultValues[envVar])
		lifetime, err := parseLifetime(value)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %v", envVar, err)
		}
		*lifetimePtr = lifetime
	}

	return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseLifetime(lifetime string) (int, error) {
	switch lifetime {
	case "permanent":
		return -1, nil
	case "":
		return 30, nil // default to 30 days
	default:
		return strconv.Atoi(lifetime)
	}
}
