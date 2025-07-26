package main

import (
	"fmt"
)

func main() {
	cfg := &Config{
		PullsMap: make(map[string][]PullRequest),
	}

	issues, err := fetchRecentIssues(365) //last year (get rid of this magic number)
	if err != nil {
		fmt.Printf("Error fetching issues: %v\n", err)
		return
	}
	cfg.Issues = issues

	cfg.TopRepos = topRepos(cfg.Issues, 10)
	enriched, err := RepoDetails(cfg.TopRepos)
	if err != nil {
		fmt.Printf("Error enriching repos: %v\n", err)
		return
	}
	cfg.TopRepos = enriched

	for _, r := range cfg.TopRepos {
		prs, _ := fetchRecentPulls(r.PullsURL, 5)
		cfg.PullsMap[r.PullsURL] = prs
	}

	displayTopRepos(cfg) // Show repos immediately at startup
	startRepl(cfg)
}
