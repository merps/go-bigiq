package main

import (
	"fmt"
	"gitlab.wirelessravens.org/go-bigiq"
)

func main() {
	// Connect to the BIQ-IP system. Enabled basic auth in BIQ - https://support.f5.com/csp/article/K43725273
	// TODO: how to push an F5-Auth-Token?
	f5 := bigiq.NewSession("10.0.90.253", "443", "admin", "SuperSecretPAssword!", nil)

	// Get a list of all VLAN's, and print their names to the console.
	//devices, err := f5.GetDevices()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(devices)
	//}

	devices, err := f5.GetDevices()
	if err != nil {
		fmt.Println(err)
	}

	for _, devices := range devices {
		fmt.Println(devices)
	}

	logsrvs, err := f5.Syslogs()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(logsrvs)
	}
	// Create a VLAN
	//f5.CreateVlan("vlan1138", 1138)
	//f5.CreateVlan("vlan421", 421)

	// Add an untagged interface to a VLAN.
	//f5.AddInterfaceToVlan("vlan1138", "1.2", false)

	// Delete a VLAN.
	//f5.DeleteVlan("vlan1138")

	// Create a couple of nodes.
	//f5.CreateNode("web-server-1", "192.168.1.50")
	//f5.CreateNode("web-server-2", "192.168.1.51")
	//f5.CreateNode("ssl-web-server-1", "10.2.2.50")
	//f5.CreateNode("ssl-web-server-2", "10.2.2.51")

	// Create a pool, and add members to it. When adding a member, you must
	// specify the port in the format of <node name>:<port>.
	//f5.CreatePool("web_farm_80_pool")
	//f5.AddPoolMember("web_farm_80_pool", "web-server-1:80")
	//f5.AddPoolMember("web_farm_80_pool", "web-server-2:80")
	//f5.CreatePool("ssl_443_pool")
	//f5.AddPoolMember("ssl_443_pool", "ssl-web-server-1:443")
	//f5.AddPoolMember("ssl_443_pool", "ssl-web-server-2:443")

	// Create a monitor, and assign it to a pool.
	//f5.CreateMonitor("web_http_monitor", "http", 5, 16, "GET /\r\n", "200 OK", "http")
	//f5.AddMonitorToPool("web_http_monitor", "web_farm_80_pool")

	// Create a virtual server, with the above pool. The third field is the subnet
	// mask, and that can either be in CIDR notation or decimal. For any/all destinations
	// and ports, use '0' for the mask and/or port.
	//f5.CreateVirtualServer("web_farm_VS", "0.0.0.0", "0.0.0.0", "web_farm_80_pool", 80)
	//f5.CreateVirtualServer("ssl_web_farm_VS", "10.1.1.0", "24", "ssl_443_pool", 443)

	// Remove a pool member.
	//f5.DeletePoolMember("web_farm_80_pool", "web-server-2:80")

	// Create a trunk, with LACP enabled.
	//f5.CreateTrunk("Aggregated", "1.2, 1.4, 1.6", true)

	// Disable a virtual address.
	//f5.VirtualAddressStatus("web_farm_VS", "disable")

	// Disable a pool member.
	//f5.PoolMemberStatus("ssl_443_pool", "ssl-web-server-1:443", "disable")

	// Create a self IP.
	//f5.CreateSelfIP("vlan1138", "10.10.10.1/24", "vlan1138")
	//f5.CreateSelfIP("vlan421", "10.10.20.1/25", "vlan421")

	// Add a static route.
	//f5.CreateRoute("servers", "10.1.1.0/24", "10.50.1.5")

	// Create a route domain.
	//f5.CreateRouteDomain("vlans", 10, true, "vlan1138, vlan421")
}
