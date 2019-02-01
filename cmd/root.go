// Copyright © 2019 Alexander Pinnecke <alexander.pinnecke@googlemail.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"
)

var logger *logrus.Logger

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gitlab-razer-device-build-gitlab",
	Short: "Watch a list of GitLab repositories and watch for failed builds using razer devices as feedback device",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger = logrus.New()

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			logger.Fatalf("failed to get verbose parameter: %v", err)
		}

		if verbose {
			logger.SetLevel(logrus.DebugLevel)
		}

		plain, err := cmd.Flags().GetBool("plain")
		if err != nil {
			logger.Fatalf("failed to get plain parameter: %v", err)
		}

		if plain {
			logger.SetFormatter(&logrus.TextFormatter{
				DisableColors:    false,
				ForceColors:      true,
				QuoteEmptyFields: true,
			})
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Run the command in verbose / debug mode")
	RootCmd.PersistentFlags().BoolP("plain", "p", false, "Run the command with plain output mode mode")

}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
