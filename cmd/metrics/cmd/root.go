package cmd

import (
	"github.com/arjungandhi/metrics/pkg/metrics"
	"github.com/arjungandhi/metrics/pkg/store"
	"github.com/spf13/cobra"
)

var s store.Store

var rootCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Track metrics from the command line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		client, err := metrics.New(nil)
		if err != nil {
			return err
		}
		s = client.Store
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
