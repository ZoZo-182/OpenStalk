package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *Config) {
	reader := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Print("OpenStalk> ")
		reader.Scan()
		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		cmd, ok := commands[words[0]]
		if !ok {
			fmt.Printf("Unknown command: %s\n", words[0])
			continue
		}
		if err := cmd.callback(cfg, words[1:]...); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

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
