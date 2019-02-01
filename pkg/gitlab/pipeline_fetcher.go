package gitlab

import (
	"fmt"

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
}

// NewPipelineFetcher returns a new PipelineFetcher instance
func NewPipelineFetcher(logger *logrus.Entry, client pipelinesClient) (*PipelineFetcher, error) {
	return &PipelineFetcher{
		logger: logger,
		client: client,
	}, nil
}

// GetPipelineStatusForProject returns the latest state for a projects pipeline, limited to given username
func (p *PipelineFetcher) GetPipelineStatusForProject(username string, projectID int) (string, error) {
	opt := &gitlab.ListProjectPipelinesOptions{
		Username: gitlab.String(username),
	}

	pipelines, _, err := p.client.ListProjectPipelines(projectID, opt)
	if err != nil {
		return "", fmt.Errorf("failed to fetch currently authenticated pipeline: %v", err)
	}

	if len(pipelines) == 0 {
		return StatusSuccess, nil
	}

	var latestStatus string
	for _, pipeline := range pipelines {
		if pipeline.Status != StatusSuccess && pipeline.Status != StatusFailed {
			p.logger.Debugf("Skipping pipeline %d due to irrelevant status %q", pipeline.Status)
			continue
		}

		p.logger.Debugf("Processing pipeline project=%d pipeline=%d state=%s ref=%s", projectID, pipeline.ID, pipeline.Status, pipeline.Ref)
		latestStatus = pipeline.Status
	}

	return latestStatus, nil
}
