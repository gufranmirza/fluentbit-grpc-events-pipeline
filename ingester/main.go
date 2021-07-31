package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/ingester/server"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Ingester",
	Short: "Ingester implementation for PLogger",
	Long:  `Starts a http server and serves the configured api`,
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "rpc-server",
	Short: "start rpc server with configured api",
	Long:  `Starts a rpc server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		fstatus, _ := cmd.Flags().GetBool("decrypt")
		server := server.NewServer(&server.Config{Decrypt: fstatus})
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().Bool("decrypt", false, "Decrypt messages received from fluentbit-agent")

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
