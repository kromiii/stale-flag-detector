package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kromiii/stale-flag-detector/config"
	"github.com/kromiii/stale-flag-detector/unleash"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	excludePotentiallyStaleFlags := flag.Bool("exclude-potentially-stale-flags", false, "Include potentially stale flags")
	flag.Parse()

	client := unleash.NewClient(cfg.UnleashAPIEndpoint, cfg.UnleashAPIToken, cfg.ProjectID, cfg)
	staleFlags, err := client.GetStaleFlags(*excludePotentiallyStaleFlags)
	if err != nil {
		fmt.Printf("Error getting stale flags: %v\n", err)
		os.Exit(1)
	}

	// Print stale flags
	fmt.Println("Stale flags:")
	for _, flag := range staleFlags {
		fmt.Printf("- %s\n", flag)
	}
}
