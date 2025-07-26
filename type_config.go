package main

type Config struct {
	Issues   []GithubIssue
	TopRepos []RepoInfo
	PullsMap map[string][]PullRequest
}
