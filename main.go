package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kromiii/stale-flag-detector/config"
	"github.com/kromiii/stale-flag-detector/unleash"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	excludePotentiallyStaleFlags := flag.Bool("exclude-potentially-stale-flags", false, "Exclude potentially stale flags")
	outputRegex := flag.Bool("output-regex", false, "Output flags as a grep-compatible regex")
	flag.Parse()

	client := unleash.NewClient(cfg.UnleashAPIEndpoint, cfg.UnleashAPIToken, cfg.ProjectID, cfg)
	staleFlags, err := client.GetStaleFlags(*excludePotentiallyStaleFlags)
	if err != nil {
		fmt.Printf("Error getting stale flags: %v\n", err)
		os.Exit(1)
	}

	if *outputRegex {
		regex := strings.Join(staleFlags, "|")
		fmt.Printf("(%s)\n", regex)
	} else {
		fmt.Println("Stale flags:")
		for _, flag := range staleFlags {
			fmt.Printf("- %s\n", flag)
		}
	}
}
