package monitor

import "github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/gitlab"

type userFetcher interface {
	GetCurrentUserName() (string, error)
}

type repoFetcher interface {
	GetProjectsWithAtLeastDevAccess() ([]gitlab.Repo, error)
}

type pipelineFetcher interface {
	GetPipelineStatusForProject(username string, projectID int) (string, error)
}

type RepoPipelineStatus struct {
	Repo   gitlab.Repo
	Status string
}
