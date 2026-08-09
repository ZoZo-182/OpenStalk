package main

import "testing"

func TestBuildSearchQueryIncludeRequiredFilters(t *testing.T) {
	query := buildSearchQuery(SearchOptions{}, "2026-06-27")

	want := "type:pr state:open created:>=2026-06-27"

	if query != want {
		t.Fatalf("expected %q, got %q", want, query)
	}
}

func TestBuildSearchQueryIncludeLanguage(t *testing.T) {
	query := buildSearchQuery(SearchOptions{
		Language: "go",
	}, "2026-06-27")

	want := "type:pr state:open created:>=2026-06-27 language:go"

	if query != want {
		t.Fatalf("expected %q, got %q", want, query)
	}
}

func TestDefaultSearchOptions(t *testing.T) {
	options := defaultSearchOptions(SearchOptions{})

	if options.Days != 1 {
		t.Fatalf("expected Days to default to 1, got %d", options.Days)
	}
	if options.Limit != 10 {
		t.Fatalf("expected Limit to default to 10, got %d", options.Limit)
	}
}
