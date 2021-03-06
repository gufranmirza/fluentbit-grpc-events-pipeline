package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/ingester/server"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/pkg/jwtauth"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Ingester",
	Short: "Ingester implementation for FB-Agent",
	Long:  `Starts a http server and serves the configured api`,
}

var (
	expiry int64
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "api",
	Short: "Start Ingester API Server",
	Long:  `Start Ingester API Server`,
	Run: func(cmd *cobra.Command, args []string) {
		fstatus, _ := cmd.Flags().GetBool("print-events")
		server := server.NewServer(&server.Config{Print: fstatus})
		server.Start()
	},
}

var authCmd = &cobra.Command{
	Use:   "access-token",
	Short: "Generate a JWT Access Token for FB-Agent",
	Long:  `Generate a JWT Access Token for FB-Agent`,
	Run: func(cmd *cobra.Command, args []string) {
		auth := jwtauth.NewJWTAuth()
		token, err := auth.Generate(&jwtauth.Claims{}, time.Duration(expiry*int64(time.Second)))
		if err != nil {
			fmt.Printf("ERROR - %v\n", err)
		}
		fmt.Println(token)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(authCmd)

	serveCmd.PersistentFlags().Bool("print-events", false, "Print events as received from Fluentbit-Agent")
	authCmd.PersistentFlags().Int64VarP(&expiry, "expiry", "", 600, "Expiry duration of token in seconds")
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
