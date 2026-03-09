package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <metric>",
	Short: "Get data points for a metric",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := s.GetMetric(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("%s (%s)\n", m.Name, m.Unit)
		for _, dp := range m.DataPoints {
			fmt.Printf("  %s  %.2f\n", dp.Time.Format("2006-01-02 15:04"), dp.Value)
			for _, item := range dp.Items {
				fmt.Printf("    - %s: %.2f\n", item.Name, item.Value)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
