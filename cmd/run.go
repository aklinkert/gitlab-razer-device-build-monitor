// Copyright © 2019 Alexander Pinnecke <alexander.pinnecke@googlemail.com>
//

package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/apinnecke/go-exitcontext"

	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/monitor"

	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/config"
	gl "github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/gitlab"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the gitlab",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New()

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			logger.Fatalf("failed to get verbose parameter: %v", err)
		}

		if verbose {
			logger.SetLevel(logrus.DebugLevel)
		}

		cfgFilePath, err := cmd.Flags().GetString("config-file")
		if err != nil {
			logger.Fatalf("failed to get config-file parameter: %v", err)
		}

		cfg, err := config.Parse(cfgFilePath)
		if err != nil {
			logger.Fatal(err)
		}

		client := gitlab.NewClient(nil, os.Getenv("GITLAB_API_TOKEN"))

		userFetcher, err := gl.NewUserFetcher(logger.WithField("module", "user_fetcher"), client.Users)
		if err != nil {
			logger.Fatalf("Failed to create UserFetcher: %v", err)
		}

		repoFetcher, err := gl.NewRepoFetcher(logger.WithField("module", "repo_fetcher"), client.Groups, cfg)
		if err != nil {
			logger.Fatalf("Failed to create RepoFetcher: %v", err)
		}

		pipelineFetcher, err := gl.NewPipelineFetcher(logger.WithField("module", "pipeline_fetcher"), client.Pipelines)
		if err != nil {
			logger.Fatalf("Failed to create PipelineFetcher: %v", err)
		}

		mon, err := monitor.New(logger.WithField("module", "monitor"), userFetcher, repoFetcher, pipelineFetcher)
		if err != nil {
			logger.Fatalf("Failed to create Monitor: %v", err)
		}

		if err := mon.UpdateStatus(); err != nil {
			logger.Fatalf("Failed to do initial status update: %v", err)
		}

		ctx := exitcontext.New()
		mon.UpdateEvery(ctx, time.Minute)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("verbose", "v", false, "Run the command in verbose / debug mode")

	runCmd.Flags().StringP("config-file", "f", filepath.Join(os.Getenv("HOME"), ".gitlab-build-monitor.json"), "Path to the config file (default: ~~.gitlab-build-gitlab.json)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
