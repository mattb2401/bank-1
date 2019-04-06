package cmd

import (
	"errors"

	"github.com/mattb2401/bank/clientServer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the sever for the bank system",
	RunE: func(cmd *cobra.Command, args []string) error {
		if mode == "" {
			return errors.New("--mode or -m  required flag required to run the server")
		}
		clientServer.RunServer(mode)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "", "Runs server using a specific mode")
	viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))
	serverCmd.MarkFlagRequired("mode")
}
