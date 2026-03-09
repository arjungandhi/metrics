package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := s.ListMetrics()
		if err != nil {
			return err
		}

		if len(names) == 0 {
			fmt.Println("No metrics tracked yet.")
			return nil
		}

		for _, name := range names {
			fmt.Println(name)
		}
		return nil
	},
}

func init() {
	metricCmd.AddCommand(listCmd)
}
