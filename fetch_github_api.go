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

type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
}

func fetchRecentCommits(repoAPIURL string, perPage int) ([]Commit, error) {
	// Extract owner/repo from API URL
	parts := strings.Split(repoAPIURL, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid repo URL format: %s", repoAPIURL)
	}
	owner := parts[4]
	repo := parts[5]

	// Construct commits API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=%d", owner, repo, perPage)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetchRecentCommits(): %v", err)
	}
	defer resp.Body.Close()

	// GitHub commits API response structure
	var githubCommits []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubCommits); err != nil {
		return nil, err
	}

	// Convert to our Commit struct
	commits := make([]Commit, len(githubCommits))
	for i, gc := range githubCommits {
		commits[i] = Commit{
			SHA:     gc.SHA,
			Message: gc.Commit.Message,
			Author:  gc.Commit.Author.Name,
		}
	}

	return commits, nil
}

// Function to convert GitHub URL to API URL
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
		Forks       int       `json:"forks_count"`
		CreatedAt   time.Time `json:"created_at"`
		ContribURL  string    `json:"contributors_url"`
		CommitsURL  string    `json:"commits_url"`
		OpenIssues  int       `json:"open_issues_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return RepoInfo{}, err
	}

	// Get recent issue count (last 7 days)
	daysAgo := 7
	recentIssueCount, recentIssues, err := getRepoIssueCount(apiURL, daysAgo)
	if err != nil {
		// If we can't get issues, still return the repo info but with 0 issues
		recentIssueCount = 0
		recentIssues = []GithubIssue{}
	}

	// Create RepoInfo struct
	repoInfo := RepoInfo{
		APIURL:        apiURL,
		Name:          repoData.Name,
		HTMLURL:       repoData.HTMLURL,
		Language:      repoData.Language,
		PullsURL:      repoData.PullsURL,
		Contributors:  0, // We could fetch this but it's expensive
		Stars:         repoData.Stars,
		Forks:         repoData.Forks,
		RecentCommits: 0, // We could fetch this but it's expensive
		CreatedAt:     repoData.CreatedAt,
		Count:         recentIssueCount, // Recent issues count
		RecentIssues:  recentIssues,
		Score:         0, // Not applicable for bookmarks
	}

	return repoInfo, nil
}

// Database structure for storing bookmarks
type BookmarkDatabase struct {
	Bookmarks []RepoInfo `json:"bookmarks"`
}

const DATABASE_FILE = "database.json"

func loadBookmarks() []RepoInfo {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Bookmark loading panic recovered: %v\n", r)
		}
	}()

	if _, err := os.Stat(DATABASE_FILE); os.IsNotExist(err) {
		return []RepoInfo{}
	}

	data, err := os.ReadFile(DATABASE_FILE)
	if err != nil {
		fmt.Printf("Error reading database file: %v\n", err)
		return []RepoInfo{}
	}

	var db BookmarkDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		fmt.Printf("Error parsing database file: %v\n", err)
		return []RepoInfo{}
	}

	return db.Bookmarks
}

// Save bookmark to database.json file
func saveBookmark(repo RepoInfo) {
	// Load existing bookmarks
	bookmarks := loadBookmarks()

	// Check if bookmark already exists
	for _, existing := range bookmarks {
		if existing.HTMLURL == repo.HTMLURL {
			return // Already bookmarked
		}
	}

	// Add new bookmark
	bookmarks = append(bookmarks, repo)

	// Save to file
	saveBookmarksToFile(bookmarks)
}

// Remove bookmark from database.json file
func removeBookmark(repoURL string) {
	bookmarks := loadBookmarks()

	// Find and remove bookmark
	for i, bookmark := range bookmarks {
		if bookmark.HTMLURL == repoURL {
			bookmarks = append(bookmarks[:i], bookmarks[i+1:]...)
			break
		}
	}

	// Save updated list to file
	saveBookmarksToFile(bookmarks)
}

// Helper function to save bookmarks to file
func saveBookmarksToFile(bookmarks []RepoInfo) {
	db := BookmarkDatabase{
		Bookmarks: bookmarks,
	}

	// Convert to JSON
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling bookmarks: %v\n", err)
		return
	}

	// Write to file
	if err := os.WriteFile(DATABASE_FILE, data, 0644); err != nil {
		fmt.Printf("Error writing database file: %v\n", err)
		return
	}
}
