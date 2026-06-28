package main

import "testing"

func TestBuildSearchQueryIncludeRequiredFilters(t *testing.T) {
	query := buildSearchQuery(SearchOptions{
		MinStars: 100,
		MaxStars: 500,
	}, "2026-06-27")

	want := "type:pr state:open created:>=2026-06-27 stars:100..500"

	if query != want {
		t.Fatalf("expected %q, got %q", want, query)
	}
}

func TestBuildSearchQueryIncludeLanguage(t *testing.T) {
	query := buildSearchQuery(SearchOptions{
		Language: "go",
		MinStars: 100,
		MaxStars: 500,
	}, "2026-06-27")

	want := "type:pr state:open created:>=2026-06-27 stars:100..500 language:go"

	if query != want {
		t.Fatalf("expected %q, got %q", want, query)
	}
}
