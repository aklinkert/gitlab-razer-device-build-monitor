package monitor

import "github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/gitlab"

type userFetcher interface {
	GetCurrentUserName() (string, error)
}

type repoFetcher interface {
	GetProjectsWithAtLeastDevAccess() ([]gitlab.Repo, error)
}

type pipelineFetcher interface {
	GetPipelineStatusForEachRef(projectID int) (gitlab.RepoStatus, error)
}

// RepoPipelineStatus carries the current status of all refs of a repo
type RepoPipelineStatus struct {
	Repo   gitlab.Repo
	Status gitlab.RepoStatus
}

// StatusNotification is distributed carrying the infos about all failed builds and the overall status
type StatusNotification struct {
	Status          string
	FailedPipelines []StatusNotificationPipeline
}

// StatusNotificationPipeline carries the info about which pipeline in which repo failed
type StatusNotificationPipeline struct {
	RepoID   int
	RepoName string
	RepoURL  string
	Branch   string
}

// NotificationReceiver is the subscriber to StatusNotifications
type NotificationReceiver interface {
	ReceiveStatusNotification(status StatusNotification)
}
