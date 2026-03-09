package cmd

import (
	"fmt"
	"os"

	"github.com/arjungandhi/health/pkg/store"
	"github.com/spf13/cobra"
)

var s store.Store

var rootCmd = &cobra.Command{
	Use:   "health",
	Short: "Track health metrics from the command line",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	ls, err := store.NewLocalStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	s = ls
}
