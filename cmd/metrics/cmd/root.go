package cmd

import (
	"github.com/arjungandhi/metrics/pkg/metrics"
	"github.com/spf13/cobra"
)

var client *metrics.Client

func completeMetricNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names, err := client.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

var rootCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Track metrics from the command line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		client, err = metrics.New(nil)
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if client != nil {
			return client.Close()
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
