package razer

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

type Device struct {
	bus         dbus.BusObject
	id, name    string
	chromaSpecs []string
}

func NewDevice(bus dbus.BusObject) (*Device, error) {
	d := &Device{
		bus: bus,
	}

	if err := d.init(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Device) init() error {
	node, err := introspect.Call(d.bus)
	if err != nil {
		return fmt.Errorf("failed to introspect device bus object: %v", err)
	}

	for _, face := range node.Interfaces {
		if face.Name == dbusInterfaceLightningBrightness {

		}
	}

}

func (d *Device) ChromaSetStatic(red, green, blue uint8) {

}
