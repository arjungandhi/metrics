package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	day        string
	labelFlags []string
)

var addCmd = &cobra.Command{
	Use:               "add <metric> <value>",
	Short:             "Add a data point to a metric",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeMetricNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		value, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return fmt.Errorf("invalid value %q: %w", args[1], err)
		}

		var ts time.Time
		if day != "" {
			ts, err = time.ParseInLocation("2006-01-02", day, time.Local)
			if err != nil {
				return fmt.Errorf("invalid day %q (expected YYYY-MM-DD): %w", day, err)
			}
		} else {
			ts = time.Now()
		}

		labels, err := parseLabels(labelFlags)
		if err != nil {
			return err
		}

		if err := client.AddAt(name, value, ts, labels); err != nil {
			return err
		}
		fmt.Printf("Added %.2f to %s\n", value, name)
		return nil
	},
}

func parseLabels(flags []string) (map[string]string, error) {
	if len(flags) == 0 {
		return nil, nil
	}
	labels := make(map[string]string, len(flags))
	for _, f := range flags {
		k, v, ok := strings.Cut(f, "=")
		if !ok {
			return nil, fmt.Errorf("invalid label %q (expected key=value)", f)
		}
		labels[k] = v
	}
	return labels, nil
}

func init() {
	addCmd.Flags().StringVarP(&day, "day", "d", "", "date for the entry (YYYY-MM-DD, defaults to today)")
	addCmd.Flags().StringArrayVarP(&labelFlags, "label", "l", nil, "label in key=value format (repeatable)")
	rootCmd.AddCommand(addCmd)
}
