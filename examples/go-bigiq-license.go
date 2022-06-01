package main

import (
	"fmt"
	"gitlab.wirelessravens.org/go-bigiq"
)

func main() {
	// Connect to the BIQ-IP system. Enabled basic auth in BIQ - https://support.f5.com/csp/article/K43725273
	// Correct by adding unknown type var - really!?
	f5, _ := bigiq.NewTokenSession("10.0.90.254", "443", "admin", "SuperSecret", "tmos", nil)

	//Get Licenses
	licenses, err := f5.GetRegPools()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(licenses)
	}

	//Post License?
	license, _ := f5.InitialActivation("OTCCU-KKAYN-KZXPF-VXKAM-SQQZCFH")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(license)
	}
}
