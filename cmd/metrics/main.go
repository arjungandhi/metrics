package main

import (
	"os"

	"github.com/arjungandhi/metrics/cmd/metrics/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
