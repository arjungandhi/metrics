package cmd

import "github.com/spf13/cobra"

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Manage health metrics",
}

func init() {
	rootCmd.AddCommand(metricCmd)
}
