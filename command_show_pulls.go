package main

import "fmt"

func callbackShowPulls(cfg *Config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: show-pulls owner/repo")
	}
	repo := args[0]
	prs, ok := cfg.PullsMap[repo]
	if !ok || len(prs) == 0 {
		fmt.Printf("No recent PRs found for %s\n", repo)
		return nil
	}
	fmt.Printf("Recent pull requests for %s:\n", repo)
	for _, pr := range prs {
		fmt.Printf(" - #%d %s (%s)\n", pr.Number, pr.Title, pr.HTMLURL)
	}
	return nil
}
