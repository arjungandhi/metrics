package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
	"github.com/spf13/cobra"
)

var (
	unit     string
	itemName string
	day      string
)

var addCmd = &cobra.Command{
	Use:   "add <metric> <value>",
	Short: "Add a data point to a metric",
	Args:  cobra.ExactArgs(2),
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

		if itemName != "" {
			item := metric.Item{Name: itemName, Value: value}
			if err := s.AddItemToDay(name, unit, item, ts); err != nil {
				return err
			}
			fmt.Printf("Added %s (%.2f) to %s\n", itemName, value, name)
		} else {
			dp := metric.DataPoint{
				Time:  ts,
				Value: value,
			}
			if err := s.AddDataPoint(name, unit, dp); err != nil {
				return err
			}
			fmt.Printf("Added %.2f to %s\n", value, name)
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&unit, "unit", "u", "", "unit of measurement (e.g. lbs, kcal, hours)")
	addCmd.Flags().StringVarP(&itemName, "item", "i", "", "item name (accumulates into the day's data point)")
	addCmd.Flags().StringVarP(&day, "day", "d", "", "date for the entry (YYYY-MM-DD, defaults to today)")
	metricCmd.AddCommand(addCmd)
}
