package cmd

import (
	"os"
	"path/filepath"

	"github.com/arjungandhi/metrics/pkg/store"
	"github.com/spf13/cobra"
)

var s store.Store

var rootCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Track metrics from the command line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dir := os.Getenv("METRICS_DIR")
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			dir = filepath.Join(home, ".local", "share", "metrics")
		}

		ls, err := store.NewLocalStore(dir)
		if err != nil {
			return err
		}
		s = ls
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
