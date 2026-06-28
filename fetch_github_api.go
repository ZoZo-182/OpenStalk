package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"net/url"
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
	MinStars int
	MaxStars int
}

func reposFromPrs(prList []PullRequest) ([]string, error) {
	repoURLs := []string{}

	for _, pr := range prList {
		url := strings.Split(pr.HTMLURL, "/pull")[0]
		repoURLs = append(repoURLs, url)
	}

	if len(repoURLs) == 0 {
		return nil, fmt.Errorf("empty / nil repo urls from prs (reposFromPrs).")
	}

	return repoURLs, nil
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

	resp, err := http.Get(requestedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prs SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, err
	}
	return prs.Items, nil
}
