package main

import (
	"fmt"
	"strings"
)

func buildSearchQuery(opts SearchOptions, cutoff string) string {
	parts := []string{
		"type:pr",
		"state:open",
		fmt.Sprintf("created:>=%s", cutoff),
		fmt.Sprintf("stars:%d..%d", opts.MinStars, opts.MaxStars),
	}

	if opts.Language != "" {
		parts = append(parts, "language:"+opts.Language)
	}

	return strings.Join(parts, " ")
}

func defaultSearchOptions(opts SearchOptions) SearchOptions {
	if opts.Days <= 0 {
		opts.Days = 1
	}
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.MinStars < 0 {
		opts.MinStars = 0
	}
	if opts.MaxStars <= 0 {
		opts.MaxStars = 500
	}
	return opts
}
