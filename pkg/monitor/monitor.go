package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/gitlab"
)

// Monitor handles all of the GitLab polling and merging of the infos
type Monitor struct {
	logger          *logrus.Entry
	userFetcher     userFetcher
	repoFetcher     repoFetcher
	pipelineFetcher pipelineFetcher

	repos          map[int]gitlab.Repo
	pipelineStatus map[int]gitlab.RepoStatus
	username       string

	notificationReceivers []NotificationReceiver
}

// New returns a new Monitor instance
func New(logger *logrus.Entry, userFetcher userFetcher, repoFetcher repoFetcher, pipelineFetcher pipelineFetcher) (*Monitor, error) {
	return &Monitor{
		logger:          logger,
		userFetcher:     userFetcher,
		repoFetcher:     repoFetcher,
		pipelineFetcher: pipelineFetcher,

		repos:          make(map[int]gitlab.Repo),
		pipelineStatus: make(map[int]gitlab.RepoStatus),

		notificationReceivers: make([]NotificationReceiver, 0),
	}, nil
}

// RegisterNotificationReceiver adds a NotificationReceiver to be notified on a StatusNotification
func (m *Monitor) RegisterNotificationReceiver(receiver NotificationReceiver) {
	m.notificationReceivers = append(m.notificationReceivers, receiver)
}

func (m *Monitor) setRepos(repos []gitlab.Repo) {
	for _, repo := range repos {
		m.repos[repo.ID] = repo
	}
}

func (m *Monitor) setRepoStatus(repo gitlab.Repo, status gitlab.RepoStatus) {
	m.pipelineStatus[repo.ID] = status
}

func (m *Monitor) getUsername() (string, error) {
	if m.username != "" {
		return m.username, nil
	}

	m.logger.Debug("Fetching gitlab username ...")

	var err error
	m.username, err = m.userFetcher.GetCurrentUserName()
	if err != nil {
		return "", fmt.Errorf("failed to fetch username: %v", err)
	}

	m.logger.Debugf("Fetching gitlab username done. Username is %s", m.username)

	return m.username, nil
}

func (m *Monitor) getRepos() (map[int]gitlab.Repo, error) {
	if len(m.repos) > 0 {
		return m.repos, nil
	}

	m.logger.Debug("Fetching gitlab repos ...")

	repos, err := m.repoFetcher.GetProjectsWithAtLeastDevAccess()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %v", err)
	}

	m.logger.Debugf("Fetching gitlab repos done. Got %d repos.", len(repos))

	m.setRepos(repos)

	return m.repos, nil
}

// UpdateEvery takes an interval and updates the status periodically
func (m *Monitor) UpdateEvery(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	var err error

	for {
		select {
		case <-ctx.Done():
			return

		case <-t.C:
			err = m.UpdateStatus()
			if err != nil {
				m.logger.Error(err)
			}
		}
	}
}

// UpdateStatus updates the status of all the repos' pipelines
func (m *Monitor) UpdateStatus() error {
	repos, err := m.getRepos()
	if err != nil {
		return err
	}

	m.logger.Infof("Updating current status ...")

	for _, repo := range repos {
		status, err := m.pipelineFetcher.GetPipelineStatusForEachRef(repo.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch pipelines for project %q: %v", repo.Name, err)
		}

		m.setRepoStatus(repo, status)
	}

	m.logger.Infof("Updating current status done. Doing Notify.")
	m.notifyStatus()

	return nil
}

// GetCurrentStatus returns a list of all repos with their corresponding branches
func (m *Monitor) GetAllCurrentStatus() ([]RepoPipelineStatus, error) {
	var status []RepoPipelineStatus

	for repoID, latestPipelineStatus := range m.pipelineStatus {
		repo := m.getRepo(repoID)
		status = append(status, RepoPipelineStatus{
			Repo:   repo,
			Status: latestPipelineStatus,
		})
	}

	return status, nil
}

func (m *Monitor) getRepo(repoID int) gitlab.Repo {
	repo, ok := m.repos[repoID]
	if !ok {
		// I guess it's okay to throw a panic as this shouldn't / can't really happen
		// ... But does it need to get checked then ... ? @TODO check back
		panic(fmt.Errorf("failed to find repo %d in state", repoID))
	}
	return repo
}

func (m *Monitor) notifyStatus() {
	username, err := m.getUsername()
	if err != nil {
		m.logger.Error(err)
	}

	overallStatus := gitlab.StatusSuccess
	var failed []StatusNotificationPipeline

	for repoID, repoStatus := range m.pipelineStatus {
		status, failedPipelines := m.getRepoStatus(username, repoID, repoStatus)
		if status == gitlab.StatusFailed {
			overallStatus = gitlab.StatusFailed
			failed = append(failed, failedPipelines...)
		}
	}

	m.dispatch(StatusNotification{
		Status:          overallStatus,
		FailedPipelines: failed,
	})
}

func (m *Monitor) getRepoStatus(username string, repoID int, repoStatus gitlab.RepoStatus) (string, []StatusNotificationPipeline) {
	overallStatus := gitlab.StatusSuccess
	var failed []StatusNotificationPipeline

	for ref, status := range repoStatus {
		if status.Status == gitlab.StatusSuccess {
			continue
		}

		if status.Username != username {
			continue
		}

		repo := m.getRepo(repoID)
		overallStatus = gitlab.StatusFailed
		failed = append(failed, StatusNotificationPipeline{
			Branch:   ref,
			RepoID:   repoID,
			RepoName: repo.Name,
			RepoURL:  repo.FullPath,
		})
	}

	return overallStatus, failed
}

func (m *Monitor) dispatch(notification StatusNotification) {
	for _, rec := range m.notificationReceivers {
		rec.ReceiveStatusNotification(notification)
	}
}
