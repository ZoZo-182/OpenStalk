package main

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestSearchIsDefaultCommand(t *testing.T) {
	parser, err := kong.New(&CLI)
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := parser.Parse([]string{})
	if err != nil {
		t.Fatal(err)
	}

	if ctx.Command() != "search" {
		t.Fatalf("expected search command, got %q", ctx.Command())
	}
}
