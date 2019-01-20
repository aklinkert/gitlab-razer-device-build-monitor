// Copyright © 2019 Alexander Pinnecke <alexander.pinnecke@googlemail.com>
//

package cmd

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/config"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/monitor"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New()

		cfgFilePath, err := cmd.Flags().GetString("config-file")
		if err != nil {
			logger.Fatalf("failed to get config-file parameter: %v", err)
		}

		cfg, err := config.Parse(cfgFilePath)
		if err != nil {
			logger.Fatal(err)
		}

		client := gitlab.NewClient(nil, os.Getenv("GITLAB_API_TOKEN"))
		fetcher, err := monitor.NewRepoFetcher(logger.WithField("module", "repo_fetcher"), client, cfg)
		if err != nil {
			logger.Fatalf("Failed to create GitLab client: %v", err)
		}

		repos, err := fetcher.Fetch()
		if err != nil {
			logger.Fatalf("Failed to fetch repos: %v", err)
		}

		for _, r := range repos {
			logger.Infof("Found repo: %v", r)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("config-file", "f", filepath.Join(os.Getenv("HOME"), ".gitlab-build-monitor.json"), "Path to the config file (default: ~~.gitlab-build-monitor.json)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
