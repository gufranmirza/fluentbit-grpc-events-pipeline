package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/router/consumer"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Router",
	Short: "Router implementation for PLogger",
	Long:  `Starts consuming events from kafka and write to specified destination`,
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Starts consuming events from kafka and write to specified destination",
	Long:  `Starts consuming events from kafka and write to specified destination`,
	Run: func(cmd *cobra.Command, args []string) {
		fstatus, _ := cmd.Flags().GetBool("decrypt-events")
		pl, _ := cmd.Flags().GetBool("print-events")
		c := consumer.NewConsumer(&consumer.Config{Decrypt: fstatus, Print: pl})
		c.Start()
		defer c.CloseConsumer()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().Bool("decrypt-events", false, "Decrypt events received from kafka")
	serveCmd.PersistentFlags().Bool("print-events", false, "Print events as received from Kafka")
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