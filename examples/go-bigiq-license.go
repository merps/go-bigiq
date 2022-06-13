package main

import (
	"fmt"
	"gitlab.wirelessravens.org/go-bigiq"
)

func main() {
	// Connect to the BIQ-IP system. Enabled basic auth in BIQ - https://support.f5.com/csp/article/K43725273
	// Correct by adding unknown type var - really!?
	f5, _ := bigiq.NewTokenSession("10.0.90.253", "443", "admin", "zun.lull-PLEW7ar", "tmos", nil)

	//Get Licenses
	licenses, err := f5.GetRegPools()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(licenses)
	}

	//Post License?
	//response, err := f5.InitialActivation("BLORB-QLEYE-NAPDR-UWTZP-WBVHUBF", "thing", "ACTIVATING_AUTOMATIC")
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
	//} else {
	//	fmt.Println(response)
	//}

	//Poll for status TODO: does this need json marshalling?
	//resRef, err := f5.PollActivation("BLORB-QLEYE-NAPDR-UWTZP-WBVHUBF")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//TODO: need to capture and feedback EULA as patch operation,
	//reference: https://clouddocs.f5.com/products/big-iq/mgmt-api/v8.1.0/HowToSamples/bigiq_public_api_wf/t_license_initial_activation.html
	//fmt.Printf("License Status: %s\n", resRef["status"])

	// Accept EULA hack
	//fmt.Print(f5.AcceptEULA("BLORB-QLEYE-NAPDR-UWTZP-WBVHUBF"))

	// Create a regPool?
	//fmt.Println(f5.CreateRegPools("bite-me", "YA BUM!!!"))

	// Get RegPools

}
