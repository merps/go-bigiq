package main

import (
	"fmt"
	"gitlab.wirelessravens.org/go-bigiq"
)

func main() {
	// Connect to the BIQ-IP system. Enabled basic auth in BIQ - https://support.f5.com/csp/article/K43725273
	// Correct by adding unknown type var - really!?
	f5, _ := bigiq.NewTokenSession("10.0.90.254", "443", "admin", "SuperSecret", "tmos", nil)

	// Get devices listed
	devices, err := f5.GetDevices()
	if err != nil {
		fmt.Println(err)
	}

	// iteration through the devices
	for _, devices := range devices {
		fmt.Println(devices)
	}

	// Get SysLog servers
	logsrvs, err := f5.Syslogs()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(logsrvs)
	}

	// Get Interfaces
	interfaces, err := f5.Interfaces()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(interfaces)
	}

	// GetDeviceID
	fmt.Println(f5.GetDeviceId("bigip2"))
}
