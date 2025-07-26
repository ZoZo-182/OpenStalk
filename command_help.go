package main

import "fmt"

func callbackHelp(cfg *Config, args ...string) error {
	fmt.Println("Welcome to OpenStalk!")
	fmt.Println("Usage:")

	availableCommands := getCommands()
	for _, cmd := range availableCommands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println("")
	return nil
}
