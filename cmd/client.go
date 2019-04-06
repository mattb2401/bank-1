package cmd

import (
	"errors"

	"github.com/mattb2401/bank/clientServer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mode string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Runs the client for the bank system.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if mode == "" {
			return errors.New("--mode or -m  required flag required to run the client")
		}
		clientServer.RunClient(mode)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "", "Runs client using a specific mode")
	viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))
	serverCmd.MarkFlagRequired("mode")
}
