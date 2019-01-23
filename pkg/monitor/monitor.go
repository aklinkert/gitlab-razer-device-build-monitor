package monitor

import "github.com/Sirupsen/logrus"

type Monitor struct {
	logger *logrus.Entry
	userFetcher userFetcher
	repoFetcher repoFetcher
	pipelineFetcher pipelineFetcher
}

func NewMonitor(logger *logrus.Entry, userFetcher userFetcher, repoFetcher repoFetcher, pipelineFetcher pipelineFetcher) (*Monitor, error) {
	return &Monitor{
		logger: logger,
		userFetcher: userFetcher,
		repoFetcher: repoFetcher,
		pipelineFetcher: pipelineFetcher,
	}, nil
}

func (m *Monitor) Run() {
	username, err := m.userFetcher.
}
