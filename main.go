package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need at least 2 arguments. Format: ./OpenStalk [language]")
	}

	daysAgo, err := strconv.Atoi(os.Args[1])

	// get slice of recent prs
	prList, err := fetchRecentPulls(daysAgo)
	if err != nil {
		fmt.Println("error fetching recent prs (main).")
	}

	repoUrl, err := reposFromPrs(prList)
	if err != nil {
		fmt.Println("error getting repo url from pr url (main).")
	}

	for i, repo := range repoUrl {
		fmt.Printf("repo %d: %v\n", i+1, repo)
	}
}

// split the

//create topraw based off pr + star (+ language later)

// Get basic repo details (name, URL, etc.)
//	enriched, err := RepoDetails(topRaw)
//	if err != nil {
//		fmt.Printf("Error getting repo details: %v\n", err)
//	}
//
//	// def rewrite
//	for i, repo := range enriched {
//		actualCount, recentIssues, err := getRepoIssueCount(repo.APIURL, 7)
//		if err != nil {
//			fmt.Printf("Warning: Could not get issue count for %s: %v\n", repo.Name, err)
//			// Keep the original count if we can't get accurate count
//		} else {
//			enriched[i].Count = actualCount
//			enriched[i].RecentIssues = recentIssues
//		}
//		// Small delay to avoid rate limiting
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	// Take top 15 for UI
//	finalCount := 15
//	if len(filtered) > finalCount {
//		filtered = filtered[:finalCount]
//	}
//
//	launchUI(filtered)
//}
