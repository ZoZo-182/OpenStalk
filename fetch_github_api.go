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
	APIURL        string
	Name          string
	Count         int
	HTMLURL       string
	Language      string
	PullsURL      string
	Contributors  int
	Stars         int
	Forks         int
	RecentCommits int
	CreatedAt     time.Time
	Score         float64
	RecentIssues  []GithubIssue
}

type PullRequest struct {
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	Number  int    `json:"number"`
}

type RepoMetadata struct {
	HTMLURL     string    `json:"html_url"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	PullsURL    string    `json:"pulls_url"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
	ContribURL  string    `json:"contributors_url"`
	CommitsURL  string    `json:"commits_url"`
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
	// Group issues by repository URL
	repoIssues := make(map[string][]GithubIssue)
	for _, issue := range issues {
		repoAPIURL := issue.RepoURL
		repoIssues[repoAPIURL] = append(repoIssues[repoAPIURL], issue)
	}

	// Create RepoInfo with counts and actual issues
	counts := make([]RepoInfo, 0, len(repoIssues))
	for repoAPIURL, issueList := range repoIssues {
		counts = append(counts, RepoInfo{
			APIURL:       repoAPIURL,
			Count:        len(issueList),
			RecentIssues: issueList, // Store actual issues
		})
	}

	// Sort by count (most active repos first)
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	if len(counts) > topN {
		counts = counts[:topN]
	}

	return counts
}

func getRepoIssueCount(repoAPIURL string, daysAgo int) (int, []GithubIssue, error) {
	// Extract owner/repo from API URL
	parts := strings.Split(repoAPIURL, "/")
	if len(parts) < 6 {
		return 0, nil, fmt.Errorf("invalid repo URL format: %s", repoAPIURL)
	}
	owner := parts[4]
	repo := parts[5]

	cutoff := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

	// Search for issues in this specific repo
	url := fmt.Sprintf(
		"https://api.github.com/search/issues?q=repo:%s/%s+type:issue+state:open+created:>%s&sort=created&order=desc&per_page=100",
		owner, repo, cutoff,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, fmt.Errorf("getRepoIssueCount(): %v", err)
	}
	defer resp.Body.Close()

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return 0, nil, err
	}

	return len(searchResponse.Items), searchResponse.Items, nil
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

func RepoDetailsEnhanced(repos []RepoInfo, daysAgo int) ([]RepoInfo, error) {
	client := http.DefaultClient

	for i, r := range repos {
		// Get basic repo metadata
		req, _ := http.NewRequest("GET", r.APIURL, nil)
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("RepoDetailsEnhanced(): %v", err)
		}
		defer resp.Body.Close()

		var repoData struct {
			HTMLURL     string `json:"html_url"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Language    string `json:"language"`
			PullsURL    string `json:"pulls_url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
			return nil, err
		}

		// Update basic fields
		repos[i].Name = repoData.Name
		repos[i].HTMLURL = repoData.HTMLURL
		repos[i].Language = repoData.Language
		repos[i].PullsURL = repoData.PullsURL

		// Get ACTUAL recent issue count for this repo
		actualCount, recentIssues, err := getRepoIssueCount(r.APIURL, daysAgo)
		if err == nil {
			repos[i].Count = actualCount // Replace with accurate count
			repos[i].RecentIssues = recentIssues
		} else {
			fmt.Printf("Warning: Could not get issue count for %s: %v\n", repos[i].Name, err)
		}

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	return repos, nil
}
