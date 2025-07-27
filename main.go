package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	// Get recent issues for discovery
	issues, err := fetchRecentIssues(7)
	if err != nil {
		fmt.Printf("Error fetching issues: %v\n", err)
		return
	}

	// Get top repos by sample frequency
	topRaw := topRepos(issues, 20)

	// Get basic repo details (name, URL, etc.)
	enriched, err := RepoDetails(topRaw)
	if err != nil {
		fmt.Printf("Error getting repo details: %v\n", err)
		return
	}

	for i, repo := range enriched {
		actualCount, recentIssues, err := getRepoIssueCount(repo.APIURL, 7)
		if err != nil {
			fmt.Printf("Warning: Could not get issue count for %s: %v\n", repo.Name, err)
			// Keep the original count if we can't get accurate count
		} else {
			enriched[i].Count = actualCount
			enriched[i].RecentIssues = recentIssues
		}
		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	// Filter out repos with 0 issues and sort by count
	var filtered []RepoInfo
	for _, repo := range enriched {
		if repo.Count > 0 {
			filtered = append(filtered, repo)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Count > filtered[j].Count
	})

	// Take top 15 for UI
	finalCount := 15
	if len(filtered) > finalCount {
		filtered = filtered[:finalCount]
	}

	launchUI(filtered)
}
