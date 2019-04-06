package cmd

import (
	"log"

	"github.com/mattb2401/bank/httpHandlers"
	"github.com/spf13/cobra"
)

// httpCmd represents the client command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Runs the http server of the client",
	Run: func(cmd *cobra.Command, args []string) {
		err := httpHandlers.RunHttpServer()
		if err != nil {
			log.Fatalf("Couldn't start HTTP server." + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
