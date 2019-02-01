package main

import (
	"encoding/json"
	"fmt"

	"github.com/godbus/dbus/introspect"

	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(fmt.Errorf("failed to connect to dbus: %v", err))
	}

	bus := conn.Object("org.razer", "/org/razer")
	describe(bus)

	var devices []string
	err = bus.Call("razer.devices.getDevices", 0).Store(&devices)
	if err != nil {
		panic(err)
	}

	fmt.Println("Available devices:")
	for _, v := range devices {
		fmt.Println(v)
	}

	if len(devices) == 0 {
		return
	}

	deviceID := devices[0]
	deviceBus := conn.Object("org.razer", dbus.ObjectPath(fmt.Sprintf("/org/razer/device/%s", deviceID)))
	describe(deviceBus)

	var deviceName string
	err = deviceBus.Call("razer.device.misc.getDeviceName", 0).Store(&deviceName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Device %s has name %q\n", deviceID, deviceName)

	deviceBus.Call("razer.device.lighting.chroma.setStatic", 0, 0, 0, 255)
	//deviceBus.Call("razer.device.lighting.chroma.setBreathSingle", 0, 0, 255, 0)
	//deviceBus.Call("razer.device.lighting.chroma.setStarlightRandom", 0, 10)
	//deviceBus.Call("razer.device.lighting.chroma.setWave", 0, 10)

}

func describe(bus dbus.BusObject) {
	node, err := introspect.Call(bus)
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(node)
	fmt.Println("")
	fmt.Println("############################")
	fmt.Println(string(data))
	fmt.Println("############################")
	fmt.Println("")
}
