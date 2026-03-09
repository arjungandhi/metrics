package cmd

import (
	"fmt"

	"github.com/arjungandhi/metrics/pkg/provider"
	"github.com/arjungandhi/metrics/providers/strava"
	"github.com/spf13/cobra"
)

// registry of all available providers.
var providers = map[string]provider.Provider{
	"strava": strava.New(),
}

func lookupProvider(name string) (provider.Provider, error) {
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q", name)
	}
	return p, nil
}

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage metric providers",
}

var providerSetupCmd = &cobra.Command{
	Use:   "setup <provider>",
	Short: "Configure a provider (OAuth, API keys, etc.)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := lookupProvider(args[0])
		if err != nil {
			return err
		}
		return p.Setup(client.Store)
	},
}

var providerSyncCmd = &cobra.Command{
	Use:   "sync [provider]",
	Short: "Sync metrics from a provider (or all providers)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			p, err := lookupProvider(args[0])
			if err != nil {
				return err
			}
			return p.Sync(client.Store)
		}

		// Sync all providers.
		var errs []error
		for name, p := range providers {
			fmt.Printf("Syncing %s...\n", name)
			if err := p.Sync(client.Store); err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", name, err))
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("sync errors: %v", errs)
		}
		return nil
	},
}

var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available providers",
	Run: func(cmd *cobra.Command, args []string) {
		for name := range providers {
			fmt.Println(name)
		}
	},
}

func init() {
	providerCmd.AddCommand(providerSetupCmd)
	providerCmd.AddCommand(providerSyncCmd)
	providerCmd.AddCommand(providerListCmd)
	rootCmd.AddCommand(providerCmd)
}
