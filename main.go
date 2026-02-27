package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {

	//create topraw based off pr + star (+ language later)

	// Get basic repo details (name, URL, etc.)
	enriched, err := RepoDetails(topRaw)
	if err != nil {
		fmt.Printf("Error getting repo details: %v\n", err)
		return
	}

	// def rewrite
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

	// Take top 15 for UI
	finalCount := 15
	if len(filtered) > finalCount {
		filtered = filtered[:finalCount]
	}

	launchUI(filtered)
}
