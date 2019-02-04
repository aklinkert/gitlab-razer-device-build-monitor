package razer

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
)

type DeviceManager struct {
	logger   *logrus.Entry
	dbus     *dbus.Conn
	razerBus dbus.BusObject
}

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

	d.razerBus = conn.Object(dbusDest, dbusRootPath)

	return nil
}

func (d *DeviceManager) getDevices() ([]string, error) {
	var devices []string
	err := d.razerBus.Call("razer.devices.getDevices", 0).Store(&devices)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get devices from razer dbus: %v", err)
	}

	return devices, nil
}

func (d *DeviceManager) getDeviceBus(deviceID string) dbus.BusObject {
	return d.dbus.Object(dbusDest, dbus.ObjectPath(fmt.Sprintf(dbusDevicePath, deviceID)))
}

func (r *DeviceManager) send

func (d *DeviceManager) sendToAll(status string) error {
	devices, err := d.getDevices()
	if err != nil {
		return err
	}

	for _, d := range devices {

	}
}
