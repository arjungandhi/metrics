package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arjungandhi/metrics/pkg/store"
	"github.com/spf13/cobra"
)

var (
	s        store.Store
	userFlag string
)

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

		// User subcommands don't require an active user.
		if cmd.Parent() == userCmd || cmd == userCmd {
			return nil
		}

		username := userFlag
		if username == "" {
			username, err = s.DefaultUser()
			if err != nil {
				return fmt.Errorf("%w\nRun 'metrics user add <name>' to create a user", err)
			}
		}

		return s.SetUser(username)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&userFlag, "user", "", "user profile to use (overrides default)")
}
