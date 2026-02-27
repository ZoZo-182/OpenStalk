package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help":       {"help", "Displays this help message", callbackHelp},
		"clear":      {"clear", "Clears the screen", callbackClear},
		"show-pulls": {"show-pulls", "Show recent pull requests for a repo", callbackShowPulls},
		"exit":       {"exit", "Exit OpenStalk CLI", callbackExit},
	}
}
