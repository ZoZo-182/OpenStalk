package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need at least 3 arguments. Format: ./OpenStalk [language] [within (days)]")
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
