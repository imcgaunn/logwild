package main

import (
	"os"

	"mcgaunn.com/logwild/pkg/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
