package gitlab

import (
	"fmt"

	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/config"

	"github.com/Sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type pipelinesClient interface {
	GetPipeline(pid interface{}, pipeline int, options ...gitlab.OptionFunc) (*gitlab.Pipeline, *gitlab.Response, error)
	ListProjectPipelines(pid interface{}, opt *gitlab.ListProjectPipelinesOptions, options ...gitlab.OptionFunc) (gitlab.PipelineList, *gitlab.Response, error)
}

// PipelineFetcher fetches pipelines and its states from the GitLab API
type PipelineFetcher struct {
	logger *logrus.Entry
	client pipelinesClient
	config *config.Config
}

// NewPipelineFetcher returns a new PipelineFetcher instance
func NewPipelineFetcher(logger *logrus.Entry, client pipelinesClient, config *config.Config) (*PipelineFetcher, error) {
	return &PipelineFetcher{
		logger: logger,
		client: client,
		config: config,
	}, nil
}

// GetPipelineStatusForProject returns the latest state for a projects pipeline
func (p *PipelineFetcher) GetPipelineStatusForEachRef(projectID int) (RepoStatus, error) {
	logger := p.logger.WithField("project", projectID)

	opt := &gitlab.ListProjectPipelinesOptions{
		Scope: gitlab.String("branches"),
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
		//Sort: gitlab.String("asc"),
	}

	pipelines, _, err := p.client.ListProjectPipelines(projectID, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currently authenticated pipeline: %v", err)
	}

	if len(pipelines) == 0 {
		return nil, nil
	}

	refStatus := make(RepoStatus)
	for _, pipeline := range pipelines {
		pipelineLogger := logger.WithFields(logrus.Fields{
			"pipeline": pipeline.ID,
			"status":   pipeline.Status,
			"ref":      pipeline.Ref,
		})

		if _, ok := refStatus[pipeline.Ref]; ok {
			pipelineLogger.Debugf("Skipping pipeline, ref is already available")
			continue
		}

		if pipeline.Status != StatusSuccess && pipeline.Status != StatusFailed {
			pipelineLogger.Debugf("Skipping pipeline due to irrelevant status %q", pipeline.Status)
			continue
		}

		pipelineLogger.Debug("Processing pipeline")

		pipelineDetails, _, err := p.client.GetPipeline(projectID, pipeline.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch details of pipeline %d (project %d): %v", pipeline.ID, projectID, err)
		}

		refStatus[pipeline.Ref] = PipelineStatus{
			ID:       pipeline.ID,
			SHA:      pipeline.SHA,
			Ref:      pipeline.Ref,
			Status:   pipeline.Status,
			Username: pipelineDetails.User.Username,
			UserID:   pipelineDetails.User.ID,
		}
	}

	return refStatus, nil
}
