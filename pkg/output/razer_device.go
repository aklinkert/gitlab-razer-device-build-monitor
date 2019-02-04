package output

import (
	"github.com/Sirupsen/logrus"
	"github.com/apinnecke/gitlab-razer-device-build-monitor/pkg/monitor"
	"github.com/godbus/dbus"
)

type RazerDeviceOutput struct {
	logger   *logrus.Entry
	dbus     *dbus.Conn
	razerBus dbus.BusObject
}

func NewRazerDeviceOutput(logger *logrus.Entry) (*RazerDeviceOutput, error) {
	o := &RazerDeviceOutput{
		logger: logger,
	}

	if err := o.init(); err != nil {
		return nil, err
	}

	return o, nil
}

// ReceiveStatusNotification changes the color of razer devices based on overall build state
func (r *RazerDeviceOutput) ReceiveStatusNotification(notification monitor.StatusNotification) {

}
