package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <metric>",
	Short: "Get data points for a metric",
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		names, err := s.ListMetrics()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := s.GetMetric(args[0])
		if err != nil {
			return err
		}

		fmt.Println(m.Name)
		for _, dp := range m.DataPoints {
			fmt.Printf("  %s  %.2f", dp.Time.Format("2006-01-02 15:04"), dp.Value)
			for k, v := range dp.Labels {
				fmt.Printf("  %s=%s", k, v)
			}
			fmt.Println()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
