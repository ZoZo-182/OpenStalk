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

// default daysAgo to 1
func fetchRecentPulls(daysAgo int, language string) ([]PullRequest, error) {
	cutoff := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

	url := fmt.Sprintf(
		"https://api.github.com/search/issues?q=type:pr+state:open+created:>=%s+stars:100..500+language:%s&sort=created&order=desc",
		cutoff, language,
	)

	//// Trim the template suffix
	//url := strings.Split(pullsURL, "{")[0] + fmt.Sprintf("?state=open&sort=created&direction=desc&per_page=%d", perPage)
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
