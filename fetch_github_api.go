package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type PullRequest struct {
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
}

type SearchResponse struct {
	Items []PullRequest `json:"items"`
}

type SearchOptions struct {
	Days     int
	Language string
	Limit    int
}

func fetchRecentPulls(opts SearchOptions) ([]PullRequest, error) {
	opts = defaultSearchOptions(opts)

	cutoff := time.Now().AddDate(0, 0, -opts.Days).Format("2006-01-02")

	query := buildSearchQuery(opts, cutoff)

	requestedURL := fmt.Sprintf(
		"https://api.github.com/search/issues?q=%s&sort=created&order=desc&per_page=%d",
		url.QueryEscape(query),
		opts.Limit,
	)

	req, err := http.NewRequest(http.MethodGet, requestedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "OpenStalk")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("Github API rate limit reached; set GITHUB_TOKEN for a higher limit")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Github API request failed: %s", resp.Status)
	}

	var prs SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, err
	}
	return prs.Items, nil
}
