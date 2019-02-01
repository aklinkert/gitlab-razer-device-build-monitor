package output

import (
	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/monitor"
)

// LogOutput handles logging of StatusNotifications
type LogOutput struct {
	logger *logrus.Entry
}

// NewLogOutput returns a new LogOutput instance
func NewLogOutout(logger *logrus.Entry) (*LogOutput, error) {
	return &LogOutput{
		logger: logger,
	}, nil
}

// ReceiveStatusNotification renders the current StatusNotification to stdout
func (l *LogOutput) ReceiveStatusNotification(notification monitor.StatusNotification) {
	l.logger.Infof("overall status: %v", notification.Status)

	for _, f := range notification.FailedPipelines {
		l.logger.Infof("> https://gitlab.com/%s/pipelines?scope=branches&page=1 @%s", f.RepoURL, f.Branch)
	}
}
