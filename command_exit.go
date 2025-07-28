package main

import (
	"fmt"
	"os"
)

func callbackExit(cfg *Config, args ...string) error {
	fmt.Println("Leaving OpenStalk - Bye!")

	os.Exit(0)
	return nil
}
