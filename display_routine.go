package main

import (
	"fmt"
	"github.com/savioxavier/termlink"
)

func displayTopRepos(cfg *Config) {
	fmt.Println("Top repositories based on recent issue activity:")
	for _, r := range cfg.TopRepos {
		link := termlink.Link(r.Name, r.HTMLURL)
		fmt.Printf("%s — %d recent issues — language: %s\n",
			link, r.Count, r.Language)
	}
	fmt.Println()
}
