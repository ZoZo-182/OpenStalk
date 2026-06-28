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
