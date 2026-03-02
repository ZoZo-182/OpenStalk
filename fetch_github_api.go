package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RepoInfo struct {
	APIURL   string
	Name     string
	HTMLURL  string
	Language string
	Stars    int
}

type PullRequest struct {
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
}

type SearchResponse struct {
	Items []PullRequest `json:"items"`
}

type RepoMetadata struct {
	HTMLURL     string `json:"html_url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
}

// wtf where are am I getting a slice of repos from without any info in them??
// I put in the info at the end so this should actually have no parameters
// and instead declare an instance of RepoInfo and get rid of the loop.
// I should turn this into a method to use on a Config struct and then use a
// loop to access the slice of repose the cfg will hold.
// that means i only really need two functions maybe.
// get the list (slice) of repos with recent prs AND whose star count is 100-1k
// this means i have to add a query about the prs and stars. store the slice in
// an instance of Config.
// and then use the repoDetails function on this instance to actuall see info
// about the repos with recent prs + certain star count
// okay no method b/c only one member, really? get the queried data, return slice,
// can use a loop on each repo to get the info, duh

func reposFromPrs(prList []PullRequest) ([]string, error) {
	repoURLs := []string{}

	for _, pr := range prList {
		url := strings.Split(pr.HTMLURL, "/pull")[0]
		repoURLs = append(repoURLs, url)
	}

	if len(repoURLs) == 0 {
		return nil, fmt.Errorf("empty / nil repo urls from prs (reposFromPrs).")
	}

	return repoURLs, nil
}

//func RepoDetails(prList []PullRequest) ([]RepoInfo, error) {
//	client := http.DefaultClient
//	for i, pr := range prList {
//		req, _ := http.NewRequest("GET", r.APIURL, nil)
//		resp, err := client.Do(req)
//		if err != nil {
//			return nil, fmt.Errorf("RepoDetails(): %v", err)
//		}
//		defer resp.Body.Close()
//
//		var repoMetaData struct {
//			HTMLURL     string `json:"html_url"`
//			Name        string `json:"name"`
//			Description string `json:"description"`
//			Language    string `json:"language"`
//			PullsURL    string `json:"pulls_url"`
//			Stars       int    `json:"stargazers_count"`
//		}
//		if err := json.NewDecoder(resp.Body).Decode(&repoMetaData); err != nil {
//			return nil, err
//		}
//
//		repos[i].Name = repoMetaData.Name
//		repos[i].HTMLURL = repoMetaData.HTMLURL
//		repos[i].Language = repoMetaData.Language
//		repos[i].PullsURL = repoMetaData.PullsURL
//		repos[i].Stars = repoMetaData.Stars
//	}
//	return repos, nil
//}

// add based on star count too
// add language filter
// default daysAgo to 1
func fetchRecentPulls(daysAgo int) ([]PullRequest, error) {
	cutoff := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

	url := fmt.Sprintf(
		"https://api.github.com/search/issues?q=type:pr+state:open+created:>=%s&sort=created&order=desc",
		cutoff,
	)

	//// Trim the template suffix
	//url := strings.Split(pullsURL, "{")[0] + fmt.Sprintf("?state=open&sort=created&direction=desc&per_page=%d", perPage)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prs SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, err
	}
	return prs.Items, nil
}
