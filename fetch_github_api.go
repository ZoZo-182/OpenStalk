package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type SearchResponse struct {
	Items []GithubIssue `json:"items"`
}

type GithubIssue struct {
	Title       string `json:"title"`
	RepoURL     string `json:"repository_url"`
	IssueURL    string `json:"html_url"`
	IssueNumber int    `json:"number"`
}

type RepoInfo struct {
	APIURL   string
	Name     string
	Count    int
	HTMLURL  string
	Language string
	PullsURL string
}

type PullRequest struct {
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	Number  int    `json:"number"`
}

func fetchRecentIssues(daysAgo int) ([]GithubIssue, error) {
	cutoff := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")
	url := fmt.Sprintf(
		"https://api.github.com/search/issues?q=type:issue+state:open+created:>%s&sort=created&order=desc",
		cutoff,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetchRecentIssues(): %v", err)
	}
	defer resp.Body.Close()

	// decode
	var recentIssueRepos SearchResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&recentIssueRepos)
	if err != nil {
		return nil, err
	}

	return recentIssueRepos.Items, nil
}

func topRepos(issues []GithubIssue, topN int) []RepoInfo {
	freq := make(map[string]int)
	for _, issue := range issues {
		repoAPIURL := issue.RepoURL
		freq[repoAPIURL]++
	}

	counts := make([]RepoInfo, 0, len(freq))
	for repoAPIURL, cnt := range freq {
		counts = append(counts, RepoInfo{APIURL: repoAPIURL, Count: cnt})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	if len(counts) > topN {
		counts = counts[:topN]
	}

	return counts
}

func RepoDetails(repos []RepoInfo) ([]RepoInfo, error) {
	client := http.DefaultClient
	for i, r := range repos {
		req, _ := http.NewRequest("GET", r.APIURL, nil)
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("RepoDetails(): %v", err)
		}
		defer resp.Body.Close()

		var repoMetaData struct {
			HTMLURL     string `json:"html_url"`
			Name        string
			Description string
			Language    string
			PullsURL    string
		}
		if err := json.NewDecoder(resp.Body).Decode(&repoMetaData); err != nil {
			return nil, err
		}

		repos[i].Name = repoMetaData.Name
		repos[i].HTMLURL = repoMetaData.HTMLURL
		repos[i].Language = repoMetaData.Language
	}
	return repos, nil
}

func fetchRecentPulls(pullsURL string, perPage int) ([]PullRequest, error) {
	// Trim the template suffix
	url := strings.Split(pullsURL, "{")[0] + fmt.Sprintf("?state=open&sort=created&direction=desc&per_page=%d", perPage)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prs []PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, err
	}
	return prs, nil
}
