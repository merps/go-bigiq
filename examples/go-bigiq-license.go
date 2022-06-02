package main

import (
	"fmt"
	"gitlab.wirelessravens.org/go-bigiq"
)

func main() {
	// Connect to the BIQ-IP system. Enabled basic auth in BIQ - https://support.f5.com/csp/article/K43725273
	// Correct by adding unknown type var - really!?
	f5, _ := bigiq.NewTokenSession("10.0.90.253", "443", "admin", "SuperSecrect", "tmos", nil)

	//Get Licenses
	//licenses, err := f5.GetRegPools()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(licenses)
	//}

	//Post License?
	//response, err := f5.InitialActivation("OTCCU-KKAYN-KZXPF-VXKAM-SQQZCFH", "thing", "ACTIVATING_AUTOMATIC")
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println("FU to that moon!!!")
	//	return
	//}
	//fmt.Println(response)

	//Removed failed activation
	//response, err := f5.RemoveActivation("OTCCU-KKAYN-KZXPF-VXKAM-SQQZCFH")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(response)

	//Poll for status TODO: does this need json marshalling?
	resRef, err := f5.PollActivation("OTCCU-KKAYN-KZXPF-VXKAM-SQQZCFH")
	if err != nil {
		fmt.Println(err)
		return
	}
	response := resRef["status"].(string)
	fmt.Println(response)
}
