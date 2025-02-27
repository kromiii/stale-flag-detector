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
	outputFormat := flag.String("output-format", string(OutputFormatMarkdownUnorderedList), "Specifies the output format")
	flag.Parse()

	if err := validateOutputFormatFlag(*outputFormat); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	client := unleash.NewClient(cfg.UnleashAPIEndpoint, cfg.UnleashAPIToken, cfg.ProjectID, cfg)
	staleFlags, err := client.GetStaleFlags(*excludePotentiallyStaleFlags)
	if err != nil {
		fmt.Printf("Error getting stale flags: %v\n", err)
		os.Exit(1)
	}

	if len(staleFlags) == 0 {
		fmt.Println("No stale flags detected")
		return
	}

	switch *outputFormat {
	case OutputFormatMarkdownUnorderedList:
		fmt.Println("Stale flags:")
		for _, flag := range staleFlags {
			fmt.Printf("- %s\n", flag)
		}
	case OutputFormatMarkdownTaskList:
		fmt.Println("Stale flags:")
		for _, flag := range staleFlags {
			fmt.Printf("- [ ] %s\n", flag)
		}
	case OutputFormatRegex:
		regex := strings.Join(staleFlags, "|")
		fmt.Printf("(%s)\n", regex)
		// we already validated the flag format, so don't neet default clause
	}
}
