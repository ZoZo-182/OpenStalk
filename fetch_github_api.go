package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Days      int
	Language string
	Limit     int
	MinStars  int
	MaxStars  int
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
	cutoff := time.Now().AddDate(0, 0, -opts.Days).Format("2006-01-02")

	var query string
	if opts.Language != "" {
		query = fmt.Sprintf("type:pr+state:open+created:>=%s+stars:100..500+language:%s", cutoff, opts.Language)
	} else {
		query = fmt.Sprintf("type:pr+state:open+created:>=%s+stars:100..500", cutoff)
	}

	url := fmt.Sprintf(
		"https://api.github.com/search/issues?q=%s&sort=created&order=desc",
		query,
	)

	resp, err := http.Get(url)
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
