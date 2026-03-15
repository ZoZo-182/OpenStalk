package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"strings"
)

var CLI struct {
	Language string `help:"Programming Language to filter PRs by." short:"l" optional:""`
	Days     int    `help:"Number of days to look back." short:"d" default:"1" optional:""`
}

func main() {
	kong.Parse(&CLI)

	// get slice of recent prs
	prList, err := fetchRecentPulls(CLI.Days, strings.ToLower(CLI.Language))
	if err != nil {
		fmt.Println("error fetching recent prs (main).")
	}

	// change error message to be more specific (lang + time not recent enough)
	repoUrl, err := reposFromPrs(prList)
	if err != nil {
		fmt.Printf("No recent PRs to %s based repositories within the last %d day(s)\n", CLI.Language, CLI.Days)
	}

	for i, repo := range repoUrl {
		fmt.Printf("repo %d: %v\n", i+1, repo)
	}

}
