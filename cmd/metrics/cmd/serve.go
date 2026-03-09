package cmd

import (
	"fmt"

	"github.com/arjungandhi/metrics/pkg/web"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := fmt.Sprintf(":%d", servePort)
		fmt.Printf("Starting metrics UI at http://localhost%s\n", addr)
		return web.Serve(addr, s)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "port to listen on")
	rootCmd.AddCommand(serveCmd)
}
