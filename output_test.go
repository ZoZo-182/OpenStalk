package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintJSONWriteResults(t *testing.T) {
	var out bytes.Buffer

	err := printJSON(&out, []PullRequest{
		{
			Title: "Slop together the solution to non existent problem",
			HTMLURL: "https://github.com/user/repo/pull/384",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := out.String()

	if !strings.Contains(got, "Slop ") {
		t.Fatalf("expected JSON output to contain PR title, got %q", got)
	}
}
