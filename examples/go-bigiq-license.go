package main

import (
	"fmt"
)

func main() {
	// Connect to the BIQ-IP system. Enabled basic auth in BIQ - https://support.f5.com/csp/article/K43725273
	// Correct by adding unknown type var - really!?
	f5, _ := bigiq.NewTokenSession("10.0.90.254", "443", "admin", "SuperSecret", "tmos", nil)

	// Licensing Initial Activation API - 1. Start activation of a license (Automatic)
	//response, err := f5.InitialActivation("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx", "this-is-auto", "ACTIVATING_AUTOMATIC")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(response)

	// Licensing Initial Activation API - 1. Start activation of a license (Manual)
	//response, err := f5.InitialActivation("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx", "this, this is manual", "ACTIVATING_MANUAL")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(response)

	// Licensing Initial Activation API - 2. Poll to get status
	fmt.Println(f5.PollActivation("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"))
	// TODO: does this need json marshalling for output?
	//resp, err := f5.PollActivation("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//} else {
	//	output, _ := json.Marshal(resp)
	//	fmt.Println(string(output))
	//}

	// Licensing Initial Activation API - 3. Complete automatic activation by accepting the EULA
	// TODO: ugly hack - as PollStatus what's best output?
	//fmt.Print(f5.AcceptEULA("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"))

	// Licensing Initial Activation API - 4. Complete manual activation by providing license text
	// TODO: magic, all the magic!  PostReqBody
	//fmt.Print("DO MAGIC!!!!"))

	//Licensing Initial Activation API - 5. Retry Failed Activation
	//response, err := f5.RemoveActivation("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//} else {
	//	fmt.Println(response)
	//}

	//Licensing Initial Activation API - 6. Remove a failed activation
	//response, err := f5.RemoveActivation("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//} else {
	//	fmt.Println(response)
	//}

	// RegPool Examples as per documentation:
	// https://clouddocs.f5.com/products/big-iq/mgmt-api/v0.0/ApiReferences/bigiq_public_api_ref/r_license_regkey_pool.html

	// Get RegPools ID
	//fmt.Println(f5.GetRegkeyPoolId("go-biq-map"))

	// Get Pool Type
	// fmt.Println(f5.GetPoolType("go-biq-map"))

	// GET to query existing RegKey pools
	//regpools, err := f5.GetRegPools()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(regpools)
	//}

	// POST to create a new RegKey pool
	//fmt.Println(f5.CreateRegPools("go-bigiq-license.go", "go-biq-wtaf"))

	// PATCH to change the name or description of a RegKey pool
	//fmt.Println(f5.ModifyRegPool("go-biq-lic", "modify-test-hack"))

	// DELETE to remove a RegKey pool
	//fmt.Println(f5.DeleteRegPool("lab-eval"))

	// POST hack for getDossier
	// fmt.Println(f5.GetDossier("xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"))
}
