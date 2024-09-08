package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kromiii/stale-flag-detector/config"
	"github.com/kromiii/stale-flag-detector/unleash"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: unleash-checker-ai <folder>")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	onlyStaleFlag := flag.Bool("only-stale", false, "Ignore potentially stale flags")
	flag.Parse()

	client := unleash.NewClient(cfg.UnleashAPIEndpoint, cfg.UnleashAPIToken, cfg.ProjectID, cfg)
	onlyStaleFlags := *onlyStaleFlag
	staleFlags, err := client.GetStaleFlags(onlyStaleFlags)
	if err != nil {
		fmt.Printf("Error getting stale flags: %v\n", err)
		return
	}

	targetFolder := os.Args[1]
	changedFiles, removedFlags, err := finder.FindAndReplaceFlags(targetFolder, staleFlags, cfg.OpenAIAPIKey)
	if err != nil {
		fmt.Printf("Error finding affected files: %v\n", err)
		os.Exit(1)
	}

	summary := report.CreateSummary(staleFlags, removedFlags)

	// chanedFilesに変更がない場合は終了
	if len(changedFiles) == 0 {
		fmt.Println("No changes required")
		return
	}

	// Create GitHub client
	githubClient := github.NewClient(cfg.GitHubToken, cfg.GitHubOwner, cfg.GitHubRepo, cfg.GitHubBaseURL)

	ctx := context.Background()
	currentTime := time.Now().Format("20060112-150405")
	branchName := "unleash-checker-updates-" + currentTime

	// Get the default branch
	defaultBranch, err := githubClient.GetDefaultBranch(ctx)
	if err != nil {
		fmt.Printf("Error getting default branch: %v\n", err)
		os.Exit(1)
	}

	// Create branch, commit changes, and create/update pull request
	pr, err := githubClient.CreateUpdatePullRequest(ctx, branchName, "Update stale Unleash flags", "Update stale Unleash flags", summary, defaultBranch, changedFiles)
	if err != nil {
		fmt.Printf("Error creating/updating pull request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Pull request created: %s\n", pr.GetHTMLURL())
}
