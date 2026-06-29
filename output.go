package main

import (
	"fmt"
	"encoding/json"
	"os"
)

func printBanner() {
	fmt.Print(`
 ________                         _________ __         .__   __
 \_____  \ ______   ____   ____  /   _____//  |______  |  | |  | __
  /   |   \____ \_/ __ \ /    \ \_____  \   __\__  \ |  | |  |/ /
 /    |    \  |_> >  ___/|   |  \/        \|  |  / __ \|  |_|    <
 \_______  /   __/ \___  >___|  /_______  /|__| (____  /____/__|_ \
         \/|__|        \/     \/        \/           \/          \/
`)
	fmt.Println("Find active projects through recent GitHub pull requests.")
	fmt.Println()
}

func printPullRequests(prs []PullRequest) {
	for i, pr := range prs {
		fmt.Printf("%d. %s\n", i+1, pr.Title)
		fmt.Printf("    %s\n", pr.HTMLURL)
		fmt.Println()
	}
}

func printJSON(results []PullRequest) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
