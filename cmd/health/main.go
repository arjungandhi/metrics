package main

import (
	"os"

	"github.com/arjungandhi/health/cmd/health/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
