package razer

import (
	"fmt"

	"github.com/godbus/dbus/introspect"

	"github.com/godbus/dbus"
)

// Device handles a single Razer Device and manages sending chroma commands
type Device struct {
	dbus           *dbus.Conn
	bus            dbus.BusObject
	id, name       string
	introspectNode *introspect.Node
}

// NewDevice returns a new Device instance
func NewDevice(dbus *dbus.Conn, id string) (*Device, error) {
	d := &Device{
		dbus: dbus,
		id:   id,
	}

	if err := d.init(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Device) init() error {
	d.getDeviceBus(d.id)

	return d.getIntrospectionNode()
}

func (d *Device) getIntrospectionNode() error {
	var err error
	d.introspectNode, err = introspect.Call(d.bus)
	if err != nil {
		return fmt.Errorf("failed to introspect device bus object: %v", err)
	}

	return nil
}

func (d *Device) getMethodDefinition(interfaceName, methodName string) *introspect.Method {
	for faceIndex, face := range d.introspectNode.Interfaces {
		if face.Name == interfaceName {
			for methodIndex, m := range face.Methods {
				if m.Name == methodName {
					return &d.introspectNode.Interfaces[faceIndex].Methods[methodIndex]
				}
			}
		}
	}

	return nil
}

func (d *Device) getDeviceBus(deviceID string) {
	d.bus = d.dbus.Object(dbusDest, dbus.ObjectPath(fmt.Sprintf(dbusDevicePath, deviceID)))
}

func (d *Device) ChromaCommand(methodName string, params ...interface{}) error {
	fullMethodName := fmt.Sprintf("%s.%s", dbusInterfaceLightningChroma, methodName)
	methodDefinition := d.getMethodDefinition(dbusInterfaceLightningChroma, methodName)
	if methodDefinition == nil {
		return fmt.Errorf("method %s not found", methodName)
	}

	if len(methodDefinition.Args) != len(params) {
		return fmt.Errorf("method %s expects %d params, got %d", methodName, len(methodDefinition.Args), len(params))
	}

	if err := d.bus.Call(fullMethodName, 0, params...).Err; err != nil {
		return fmt.Errorf("failed to send command %s to device %s: %v", methodName, d.id, err)
	}

	return nil
}
