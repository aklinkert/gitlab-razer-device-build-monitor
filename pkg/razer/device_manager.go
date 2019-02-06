package razer

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
)

// DeviceManager handles multiple razer devices, controlling them via dbus
type DeviceManager struct {
	logger   *logrus.Entry
	dbus     *dbus.Conn
	razerBus dbus.BusObject
}

// NewDeviceManager returns a new DeviceManager instance
func NewDeviceManager(logger *logrus.Entry) (*DeviceManager, error) {
	o := &DeviceManager{
		logger: logger,
	}

	if err := o.init(); err != nil {
		return nil, err
	}

	return o, nil
}

func (d *DeviceManager) init() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to dbus: %v", err)
	}

	d.dbus = conn
	d.razerBus = conn.Object(dbusDest, dbusRootPath)

	return nil
}

func (d *DeviceManager) getDevices() ([]*Device, error) {
	var deviceIDs []string
	err := d.razerBus.Call("razer.devices.getDevices", 0).Store(&deviceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices from razer dbus: %v", err)
	}

	var devices []*Device
	for _, deviceID := range deviceIDs {
		device, err := NewDevice(d.dbus, deviceID)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (d *DeviceManager) SendChromaCommandToAll(method string, params ...interface{}) error {
	devices, err := d.getDevices()
	if err != nil {
		return err
	}

	for _, device := range devices {
		d.logger.Debugf("Calling %s on device %s with params %s", method, device.id, params)
		if err := device.ChromaCommand(method, params...); err != nil {
			return err
		}
	}

	return nil
}
