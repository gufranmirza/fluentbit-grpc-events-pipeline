package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/ingester/server"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/jwtauth"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Ingester",
	Short: "Ingester implementation for PLogger",
	Long:  `Starts a http server and serves the configured api`,
}

var (
	expiry int64
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "rpc-server",
	Short: "Start Ingester API Server",
	Long:  `Start Ingester API Server`,
	Run: func(cmd *cobra.Command, args []string) {
		fstatus, _ := cmd.Flags().GetBool("decrypt")
		server := server.NewServer(&server.Config{Decrypt: fstatus})
		server.Start()
	},
}

var authCmd = &cobra.Command{
	Use:   "access-token",
	Short: "Generate a JWT Access Token for Collector Agent",
	Long:  `Generate a JWT Access Token for Collector Agent`,
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
