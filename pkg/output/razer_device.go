package output

import (
	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/gitlab"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/monitor"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/razer"
)

type RazerDeviceOutput struct {
	logger        *logrus.Entry
	deviceManager *razer.DeviceManager
}

func NewRazerDeviceOutput(logger *logrus.Entry) (*RazerDeviceOutput, error) {
	o := &RazerDeviceOutput{
		logger: logger,
	}

	var err error
	o.deviceManager, err = razer.NewDeviceManager(logger)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// ReceiveStatusNotification changes the color of razer devices based on overall build state
func (r *RazerDeviceOutput) ReceiveStatusNotification(notification monitor.StatusNotification) {
	if notification.Status == gitlab.StatusSuccess {
		err := r.deviceManager.SendChromaCommandToAll("setStatic", 0, 0, 255)
		if err != nil {
			r.logger.Error(err)
		}
	} else {
		err := r.deviceManager.SendChromaCommandToAll("setBreathSingle", 255, 0, 0)
		if err != nil {
			r.logger.Error(err)
		}
	}
}
