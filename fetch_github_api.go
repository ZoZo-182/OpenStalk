package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)


type RepoInfo struct {
	APIURL        string
	Name          string
	Count         int
	HTMLURL       string
	Language      string
	PullsURL      string
	Stars         int
	Forks         int
	CreatedAt     time.Time
	Score         float64
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
	CreatedAt   time.Time `json:"created_at"`
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
			Name        string `json:"name"`
			Description string `json:"description"`
			Language    string `json:"language"`
			PullsURL    string `json:"pulls_url"` // Add this line
		}
		if err := json.NewDecoder(resp.Body).Decode(&repoMetaData); err != nil {
			return nil, err
		}

		repos[i].Name = repoMetaData.Name
		repos[i].HTMLURL = repoMetaData.HTMLURL
		repos[i].Language = repoMetaData.Language
		repos[i].PullsURL = repoMetaData.PullsURL // Add this line
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

// Function to convert GitHub URL to API URL
// is this for a search function within the tui?
func githubURLToAPI(githubURL string) (string, error) {
	// Remove trailing slash and .git if present
	githubURL = strings.TrimSuffix(strings.TrimSuffix(githubURL, "/"), ".git")

	// Extract owner and repo from GitHub URL using regex
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)
	matches := re.FindStringSubmatch(githubURL)

	if len(matches) != 3 {
		return "", fmt.Errorf("invalid GitHub URL format: %s", githubURL)
	}

	owner := matches[1]
	repo := matches[2]

	return fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), nil
}

// Function to fetch detailed repository information from GitHub URL
func fetchBookmarkRepoDetails(githubURL string) (RepoInfo, error) {
	// Convert GitHub URL to API URL
	apiURL, err := githubURLToAPI(githubURL)
	if err != nil {
		return RepoInfo{}, err
	}

	// Fetch basic repository metadata
	resp, err := http.Get(apiURL)
	if err != nil {
		return RepoInfo{}, fmt.Errorf("fetchBookmarkRepoDetails(): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RepoInfo{}, fmt.Errorf("repository not found or not accessible: %s", githubURL)
	}

	var repoData struct {
		HTMLURL     string    `json:"html_url"`
		Name        string    `json:"name"`
		FullName    string    `json:"full_name"`
		Description string    `json:"description"`
		Language    string    `json:"language"`
		PullsURL    string    `json:"pulls_url"`
		Stars       int       `json:"stargazers_count"`
		CreatedAt   time.Time `json:"created_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return RepoInfo{}, err
	}

	// Create RepoInfo struct
	repoInfo := RepoInfo{
		APIURL:        apiURL,
		Name:          repoData.Name,
		HTMLURL:       repoData.HTMLURL,
		Language:      repoData.Language,
		PullsURL:      repoData.PullsURL,
		Stars:         repoData.Stars,
		Forks:         repoData.Forks,
		CreatedAt:     repoData.CreatedAt,
	}

	return repoInfo, nil
}
