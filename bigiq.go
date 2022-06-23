package bigiq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	uriRegkey      = "regkey"
	uriLicenses    = "licenses"
	uriResolver    = "resolver"
	uriDevicegroup = "device-groups"
	uriCmBigIQ     = "cm-bigip-allBigIpDevices"
	uriDevice      = "device"
	uriMembers     = "members"
	uriTasks       = "tasks"
	uriManagement  = "member-management"
	uriPurchased   = "purchased-pool"
	uriPool        = "pool"
)

var tenantProperties []string = []string{"class", "constants", "controls", "defaultRouteDomain", "enable", "label", "optimisticLockKey", "remark"}

type BigIqDevice struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port,omitempty"`
}

var defaultConfigOptions = &ConfigOptions{
	APICallTimeout: 60 * time.Second,
}

type ConfigOptions struct {
	APICallTimeout time.Duration
}

// BigIQ is a container for our session state.
type BigIQ struct {
	Host      string
	User      string
	Password  string
	Token     string // if set, will be used instead of User/Password
	Transport *http.Transport
	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent     string
	Teem          bool
	ConfigOptions *ConfigOptions
}

// APIRequest builds our request before sending it to the server.
type APIRequest struct {
	Method      string
	URL         string
	Body        string
	ContentType string
}

// Upload contains information about a file upload status
type Upload struct {
	RemainingByteCount int64          `json:"remainingByteCount"`
	UsedChunks         map[string]int `json:"usedChunks"`
	TotalByteCount     int64          `json:"totalByteCount"`
	LocalFilePath      string         `json:"localFilePath"`
	TemporaryFilePath  string         `json:"temporaryFilePath"`
	Generation         int            `json:"generation"`
	LastUpdateMicros   int            `json:"lastUpdateMicros"`
}

// RequestError contains information about any error we get from a request.
type RequestError struct {
	Code       int      `json:"code,omitempty"`
	Message    string   `json:"message,omitempty"`
	ErrorStack []string `json:"errorStack,omitempty"`
}

// Error returns the error message.
func (r *RequestError) Error() error {
	if r.Message != "" {
		return errors.New(r.Message)
	}

	return nil
}

type DeviceRef struct {
	Link string `json:"link"`
}

type ManagedDevice struct {
	DeviceReference DeviceRef `json:"deviceReference"`
}

type UnmanagedDevice struct {
	DeviceAddress string `json:"deviceAddress"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	HTTPSPort     int    `json:"httpsPort,omitempty"`
}

type regKeyPools struct {
	//Items      []struct {
	//      ID       string `json:"id"`
	//      Name     string `json:"name"`
	//      SortName string `json:"sortName"`
	//} `json:"items"`
	RegKeyPoollist     []regKeyPool `json:"items"`
	RegKeyPoolSelfLink string       `json:"selfLink"`
}

type regKeyPool struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SelfLink string `json:"selfLink"`
	SortName string `json:"sortName"`
}

type devicesList struct {
	DevicesInfo []deviceInfo `json:"items"`
}
type deviceInfo struct {
	Address           string `json:"address"`
	DeviceURI         string `json:"deviceUri"`
	Hostname          string `json:"hostname"`
	HTTPSPort         int    `json:"httpsPort"`
	IsClustered       bool   `json:"isClustered"`
	MachineID         string `json:"machineId"`
	ManagementAddress string `json:"managementAddress"`
	McpDeviceName     string `json:"mcpDeviceName"`
	Product           string `json:"product"`
	SelfLink          string `json:"selfLink"`
	State             string `json:"state"`
	UUID              string `json:"uuid"`
	Version           string `json:"version"`
}

type MembersList struct {
	Members []memberDetail `json:"items"`
}

type memberDetail struct {
	AssignmentType  string `json:"assignmentType"`
	DeviceAddress   string `json:"deviceAddress"`
	DeviceMachineID string `json:"deviceMachineId"`
	DeviceName      string `json:"deviceName"`
	ID              string `json:"id"`
	Message         string `json:"message"`
	Status          string `json:"status"`
}

type regKeyAssignStatus struct {
	ID             string `json:"id"`
	DeviceAddress  string `json:"deviceAddress"`
	AssignmentType string `json:"assignmentType"`
	DeviceName     string `json:"deviceName"`
	Status         string `json:"status"`
}

type LicenseDetails struct {
	RegKey string `json:"regKey"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type LicenseParam struct {
	Address         string `json:"address,omitempty"`
	Port            int    `json:"port,omitempty"`
	AssignmentType  string `json:"assignmentType,omitempty"`
	Command         string `json:"command,omitempty"`
	Hypervisor      string `json:"hypervisor,omitempty"`
	LicensePoolName string `json:"licensePoolName,omitempty"`
	MacAddress      string `json:"macAddress,omitempty"`
	Password        string `json:"password,omitempty"`
	SkuKeyword1     string `json:"skuKeyword1,omitempty"`
	SkuKeyword2     string `json:"skuKeyword2,omitempty"`
	Tenant          string `json:"tenant,omitempty"`
	UnitOfMeasure   string `json:"unitOfMeasure,omitempty"`
	User            string `json:"user,omitempty"`
}

type LicenseEula struct {
	Status string `json:"status"`
	Eula   string `json:"eulaText"`
}

type RegPool struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type BigiqAs3AllTaskType struct {
	Items []BigiqAs3TaskType `json:"items,omitempty"`
}

type BigiqAs3TaskType struct {
	Code int64 `json:"code,omitempty"`
	//ID string `json:"id,omitempty"`
	//Declaration struct{} `json:"declaration,omitempty"`
	Results []BigiqResults `json:"results,omitempty"`
}
type BigiqResults struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	//      LineCount int64  `json:"lineCount,omitempty"`
	Host    string `json:"host,omitempty"`
	Tenant  string `json:"tenant,omitempty"`
	RunTime int64  `json:"runTime,omitempty"`
}

func (b *BigIQ) InitialActivation(regkey, name, status string) (string, error) {
	license := LicenseDetails{
		RegKey: regkey,
		Name:   name,
		Status: status,
	}
	licResp, err := b.postReq(license, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriInitActivation)
	respRef := make(map[string]interface{})
	_ = json.Unmarshal(licResp, &respRef)
	if err != nil {
		errMsg := respRef["message"].(string)
		return errMsg, err
	}
	statusMsg := respRef["message"].(string)
	time.Sleep(5 * time.Second)
	return statusMsg, nil
}

func (b *BigIQ) PollActivation(regkey string) (map[string]interface{}, error) {
	respRef := make(map[string]interface{})
	err, _ := b.getForEntity(&respRef, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriInitActivation, regkey)
	if err != nil {
		return nil, err
	}
	pollStatus, ok := respRef["status"]
	if ok {
		pollStatus = pollStatus.(string)
	} else {
		return nil, fmt.Errorf("license status not available")
	}
	return respRef, nil
}

// AcceptEULA TODO: add RegPool calls and what do I return? HTTPError and what else?
func (b *BigIQ) AcceptEULA(regkey string) error {
	// patchRef := make(map[string]interface{})
	respRef, err := b.PollActivation(regkey)
	if err != nil {
		return err
	} else {
		patchRef := LicenseEula{
			Status: activationAutoEULA,
			Eula:   respRef["eulaText"].(string),
		}
		eulaResp := b.patch(patchRef, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriInitActivation, regkey)
		//respRef := make(map[string]interface{})
		//_ = json.Unmarshal(eulaResp, &respRef)
		fmt.Println(eulaResp)
	}
	return nil
}

func (b *BigIQ) RemoveActivation(regkey string) (string, error) {
	licResp, err := b.deleteReq(uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriInitActivation, regkey)
	if err != nil {
		return "dragons here", err
	}
	fmt.Println(licResp)
	return "", nil
}

func (b *BigIQ) PostLicense(config *LicenseParam) (string, error) {
	log.Printf("[INFO] %v license to BigIP device:%v from BIGIQ", config.Command, config.Address)
	resp, err := b.postReq(config, uriMgmt, uriCm, uriDevice, uriTasks, uriLicensing, uriPool, uriManagement)
	if err != nil {
		return "", err
	}
	respRef := make(map[string]interface{})
	_ = json.Unmarshal(resp, &respRef)
	respID := respRef["id"].(string)
	time.Sleep(5 * time.Second)
	return respID, nil
}

func (b *BigIQ) GetLicenseStatus(id string) (map[string]interface{}, error) {
	licRes := make(map[string]interface{})
	err, _ := b.getForEntity(&licRes, uriMgmt, uriCm, uriDevice, uriTasks, uriLicensing, uriPool, uriManagement, id)
	if err != nil {
		return nil, err
	}
	licStatus, ok := licRes["status"]
	if ok {
		licStatus = licStatus.(string)
	} else {
		return nil, fmt.Errorf("license status not available")
	}
	for licStatus != "FINISHED" {
		//log.Printf(" status response is :%s", licStatus)
		if licStatus == "FAILED" {
			log.Println("[ERROR]License assign/revoke status failed")
			return licRes, nil
		}
		return b.GetLicenseStatus(id)
	}
	log.Printf("License Assignment is :%s", licStatus)
	return licRes, nil
}

func (b *BigIQ) GetDeviceLicenseStatus(path ...string) (string, error) {
	licRes := make(map[string]interface{})
	err, _ := b.getForEntity(&licRes, path...)
	if err != nil {
		return "", err
	}
	//log.Printf(" Initial status response is :%s", licRes["status"])
	return licRes["status"].(string), nil
}

// TODO: need to return json/map for details of creation
// TODO: additional - the PatchRegPools function just returns taskid
func (b *BigIQ) CreateRegPool(description, name string) (string, error) {
	// var self regKeyPool
	poolReq := RegPool{
		Description: description,
		Name:        name,
	}
	resp, err := b.postReq(poolReq, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return "", err
	}
	return string(resp), err
}

func (b *BigIQ) PatchRegPool(description, name string) error {
	poolPatch := RegPool{
		Description: description,
		Name:        name,
	}
	poolID, err := b.GetRegkeyPoolId(name)
	if err != nil {
		return nil
	} else {
		return b.patch(poolPatch, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolID)
	}
}

func (b *BigIQ) ModifyRegPool(name, description string) error {
	regkeyPool, _ := b.GetRegkeyPoolId(name)
	fmt.Println(regkeyPool)
	config := RegPool{
		Name:        name,
		Description: description,
	}
	return b.patch(config, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, regkeyPool)
}

func (b *BigIQ) DeleteRegPool(name string) error {
	poolId, err := b.GetRegkeyPoolId(name)
	if err != nil {
		return err
	} else {
		return b.delete(uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriInitActivation, poolId)
	}
}

func (b *BigIQ) GetRegPools() (*regKeyPools, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIQ) GetPoolType(poolName string) (*regKeyPool, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return nil, err
	}
	for _, pool := range self.RegKeyPoollist {
		if pool.Name == poolName {
			return &pool, nil
		}
	}
	return nil, nil
}

func (b *BigIQ) GetManagedDevices() (*devicesList, error) {
	var self devicesList
	err, _ := b.getForEntity(&self, uriMgmt, uriShared, uriResolver, uriDevicegroup, uriCmBigIQ, uriDevices)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (b *BigIQ) GetDeviceId(deviceName string) (string, error) {
	var self devicesList
	err, _ := b.getForEntity(&self, uriMgmt, uriShared, uriResolver, uriDevicegroup, uriCmBigIQ, uriDevices)
	if err != nil {
		return "", err
	}
	for _, d := range self.DevicesInfo {
		log.Printf("Address=%v,Hostname=%v,UUID=%v", d.Address, d.Hostname, d.UUID)
		if d.Address == deviceName || d.Hostname == deviceName || d.UUID == deviceName {
			log.Printf("SelfLink Type=%T,SelfLink=%v", d.SelfLink, d.SelfLink)
			return d.SelfLink, nil
		}
	}
	return "", nil
}

func (b *BigIQ) GetRegkeyPoolId(poolName string) (string, error) {
	var self regKeyPools
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses)
	if err != nil {
		return "", err
	}
	for _, pool := range self.RegKeyPoollist {
		if pool.Name == poolName {
			return pool.ID, nil
		}
	}
	return "", nil
}

func (b *BigIQ) RegkeylicenseAssign(config interface{}, poolId string, regKey string) (*memberDetail, error) {
	resp, err := b.postReq(config, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers)
	if err != nil {
		return nil, err
	}
	var resp1 regKeyAssignStatus
	err = json.Unmarshal(resp, &resp1)
	if err != nil {
		return nil, err
	}
	return b.GetMemberStatus(poolId, regKey, resp1.ID)
}

func (b *BigIQ) GetMemberStatus(poolId, regKey, memId string) (*memberDetail, error) {
	var self memberDetail
	err, _ := b.getForEntity(&self, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return nil, err
	}
	for self.Status != "LICENSED" {
		log.Printf("Member status:%+v", self.Status)
		if self.Status == "INSTALLATION_FAILED" {
			return &self, fmt.Errorf("INSTALLATION_FAILED with %s", self.Message)
		}
		return b.GetMemberStatus(poolId, regKey, memId)
	}
	return &self, nil
}
func (b *BigIQ) RegkeylicenseRevoke(poolId, regKey, memId string) error {
	log.Printf("Deleting License for Member:%+v", memId)
	_, err := b.deleteReq(uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	r1 := make(map[string]interface{})
	err, _ = b.getForEntity(&r1, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	log.Printf("Response after delete:%+v", r1)
	return nil
}
func (b *BigIQ) LicenseRevoke(config interface{}, poolId, regKey, memId string) error {
	log.Printf("Deleting License for Member:%+v from LicenseRevoke", memId)
	_, err := b.deleteReqBody(config, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	r1 := make(map[string]interface{})
	err, _ = b.getForEntity(&r1, uriMgmt, uriCm, uriDevice, uriLicensing, uriPool, uriRegkey, uriLicenses, poolId, uriOfferings, regKey, uriMembers, memId)
	if err != nil {
		return err
	}
	log.Printf("Response after delete:%+v", r1)
	return nil
}
func (b *BigIQ) PostAs3Bigiq(as3NewJson string) (error, string) {
	resp, err := b.postReq(as3NewJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare)
	if err != nil {
		return err, ""
	}
	var taskList BigiqAs3TaskType
	tenant_list, tenant_count, _ := b.GetTenantList(as3NewJson)
	json.Unmarshal(resp, &taskList)
	successfulTenants := make([]string, 0)
	if taskList.Code != 200 && taskList.Code != 0 {
		i := tenant_count - 1
		success_count := 0
		for i >= 0 {
			if taskList.Results[i].Code == 200 {
				successfulTenants = append(successfulTenants, taskList.Results[i].Tenant)
				success_count++
			}
			if taskList.Results[i].Code >= 400 {
				log.Printf("[ERROR] : HTTP %d :: %s for tenant %v", taskList.Results[i].Code, taskList.Results[i].Message, taskList.Results[i].Tenant)
			}
			i = i - 1
		}
		if success_count == tenant_count {
			log.Printf("[DEBUG]Sucessfully Created tenants  = %v", tenant_list)
		} else if success_count == 0 {
			return errors.New(fmt.Sprintf("Tenant Creation failed")), ""
		} else {
			finallist := strings.Join(successfulTenants[:], ",")
			return errors.New(fmt.Sprintf("Partial Success")), finallist
		}
	}
	return nil, tenant_list

}

func (b *BigIQ) GetAs3Bigiq(targetRef, tenantRef string) (string, error) {
	as3Json := make(map[string]interface{})
	as3Json["class"] = "AS3"
	as3Json["action"] = "deploy"
	as3Json["persist"] = true
	//var adcJson
	//adcJson := make(map[string]interface{})
	//adcJson := []map[string]interface{}{}
	var adcJson interface{}
	tenantList := strings.Split(tenantRef, ",")
	//log.Printf("[DEBUG] tenantList:%+v",tenantList)
	err, ok := b.getForEntityNew(&adcJson, uriMgmt, uriShared, uriAppsvcs, uriDeclare, tenantRef)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	as3JsonNew := make(map[string]interface{})
	as3jsonType := reflect.TypeOf(adcJson).Kind()
	//log.Printf("[DEBUG] as3jsonType:%+v",as3jsonType)
	if as3jsonType == reflect.Map {
		adcJsonvalue := adcJson.(map[string]interface{})
		if adcJsonvalue["target"].(map[string]interface{})["address"].(string) == targetRef {
			for _, name := range tenantList {
				if adcJsonvalue[name] != nil {
					for k, v := range adcJsonvalue[name].(map[string]interface{}) {
						if !contains(tenantProperties, k) {
							delete(v.(map[string]interface{}), "schemaOverlay")
							for _, v1 := range v.(map[string]interface{}) {
								if reflect.TypeOf(v1).Kind() == reflect.Map && v1.(map[string]interface{})["class"] == "Service_HTTP" {
									if _, ok := v1.(map[string]interface{})["pool"]; ok {
										ss := v1.(map[string]interface{})["pool"].(string)
										ss1 := strings.Split(ss, "/")
										v1.(map[string]interface{})["pool"] = ss1[len(ss1)-1]
									}
								}
							}
						}
					}
					as3JsonNew[name] = adcJsonvalue[name]
					//delete(adcJsonvalue[name].(map[string]interface{}),"schemaOverlay")
					as3JsonNew["id"] = adcJsonvalue["id"]
					as3JsonNew["class"] = adcJsonvalue["class"]
					as3JsonNew["label"] = adcJsonvalue["label"]
					as3JsonNew["remark"] = adcJsonvalue["remark"]
					as3JsonNew["target"] = adcJsonvalue["target"]
					//as3JsonNew["updateMode"] = adcJsonvalue["updateMode"]
					as3JsonNew["schemaVersion"] = adcJsonvalue["schemaVersion"]
				}
			}
		}
	} else {
		for _, adcJsonvalue1 := range adcJson.([]interface{}) {
			adcJsonvalue := adcJsonvalue1.(map[string]interface{})
			if adcJsonvalue["target"].(map[string]interface{})["address"].(string) == targetRef {
				for _, name := range tenantList {
					if adcJsonvalue[name] != nil {
						for k, v := range adcJsonvalue[name].(map[string]interface{}) {
							if !contains(tenantProperties, k) {
								delete(v.(map[string]interface{}), "schemaOverlay")
								for _, v1 := range v.(map[string]interface{}) {
									if reflect.TypeOf(v1).Kind() == reflect.Map && v1.(map[string]interface{})["class"] == "Service_HTTP" {
										if _, ok := v1.(map[string]interface{})["pool"]; ok {
											ss := v1.(map[string]interface{})["pool"].(string)
											ss1 := strings.Split(ss, "/")
											v1.(map[string]interface{})["pool"] = ss1[len(ss1)-1]
										}
									}
								}
								//if val, ok := v.(map[string]interface{})["serviceMain"]; ok {
								//      ss := val.(map[string]interface{})["pool"].(string)
								//      ss1 := strings.Split(ss, "/")
								//      val.(map[string]interface{})["pool"] = ss1[len(ss1)-1]
								//}
							}
						}
						as3JsonNew[name] = adcJsonvalue[name]
						//delete(adcJsonvalue[name].(map[string]interface{}),"schemaOverlay")
						as3JsonNew["id"] = adcJsonvalue["id"]
						as3JsonNew["class"] = adcJsonvalue["class"]
						as3JsonNew["label"] = adcJsonvalue["label"]
						as3JsonNew["remark"] = adcJsonvalue["remark"]
						as3JsonNew["target"] = adcJsonvalue["target"]
						//as3JsonNew["updateMode"] = adcJsonvalue["updateMode"]
						as3JsonNew["schemaVersion"] = adcJsonvalue["schemaVersion"]
					}
				}
			}
		}
	}
	as3Json["declaration"] = as3JsonNew
	out, _ := json.Marshal(as3Json)
	as3String := string(out)
	return as3String, nil
}

func (b *BigIQ) DeleteAs3Bigiq(as3NewJson string, tenantName string) (error, string) {
	as3Json, err := tenantTrimToDelete(as3NewJson)
	if err != nil {
		log.Println("[ERROR] Error in trimming the as3 json")
		return err, ""
	}
	return b.post(as3Json, uriMgmt, uriShared, uriAppsvcs, uriDeclare), ""
}

func tenantTrimToDelete(resp string) (string, error) {
	jsonRef := make(map[string]interface{})
	json.Unmarshal([]byte(resp), &jsonRef)

	for key, value := range jsonRef {
		if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
			for k, v := range rec {
				if k == "target" && reflect.ValueOf(v).Kind() == reflect.Map {
					continue
				}
				if rec2, ok := v.(map[string]interface{}); ok {
					for k1, v1 := range rec2 {
						if k1 != "class" && v1 != "Tenant" {
							delete(rec2, k1)
						}
					}

				}
			}
		}
	}

	b, err := json.Marshal(jsonRef)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

// NewSession sets up our connection to the BIG-IP system.
func NewSession(host, port, user, passwd string, configOptions *ConfigOptions) *BigIQ {
	var url string
	if !strings.HasPrefix(host, "http") {
		url = fmt.Sprintf("https://%s", host)
	} else {
		url = host
	}
	if port != "" {
		url = url + ":" + port
	}
	if configOptions == nil {
		configOptions = defaultConfigOptions
	}
	return &BigIQ{
		Host:     url,
		User:     user,
		Password: passwd,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		ConfigOptions: configOptions,
	}
}

// NewTokenSession sets up our connection to the BIG-IP system, and
// instructs the session to use token authentication instead of Basic
// Auth. This is required when using an external authentication
// provider, such as Radius or Active Directory. loginProviderName is
// probably "tmos" but your environment may vary.
func NewTokenSession(host, port, user, passwd, loginProviderName string, configOptions *ConfigOptions) (b *BigIQ, err error) {
	type authReq struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		LoginProviderName string `json:"loginProviderName"`
	}
	type authResp struct {
		Token struct {
			Token string
		}
	}

	auth := authReq{
		user,
		passwd,
		loginProviderName,
	}

	marshalJSON, err := json.Marshal(auth)
	if err != nil {
		return
	}

	req := &APIRequest{
		Method:      "post",
		URL:         "mgmt/shared/authn/login",
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	b = NewSession(host, port, user, passwd, configOptions)
	resp, err := b.APICall(req)
	if err != nil {
		return
	}

	if resp == nil {
		err = fmt.Errorf("unable to acquire authentication token")
		return
	}

	var aresp authResp
	err = json.Unmarshal(resp, &aresp)
	if err != nil {
		return
	}

	if aresp.Token.Token == "" {
		err = fmt.Errorf("unable to acquire authentication token")
		return
	}

	b.Token = aresp.Token.Token

	return
}

// APICall is used to query the BIG-IP web API.
func (b *BigIQ) APICall(options *APIRequest) ([]byte, error) {
	var req *http.Request
	client := &http.Client{
		Transport: b.Transport,
		Timeout:   b.ConfigOptions.APICallTimeout,
	}
	var format string
	if strings.Contains(options.URL, "mgmt/") {
		format = "%s/%s"
	} else {
		format = "%s/mgmt/tm/%s"
	}
	url := fmt.Sprintf(format, b.Host, options.URL)
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), url, body)
	if b.Token != "" {
		req.Header.Set("X-F5-Auth-Token", b.Token)
	} else if options.URL != "mgmt/shared/authn/login" {
		req.SetBasicAuth(b.User, b.Password)
	}

	//fmt.Println("REQ -- ", options.Method, " ", url," -- ",options.Body)

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		if res.Header["Content-Type"][0] == "application/json" {
			return data, b.checkError(data)
		}

		return data, errors.New(fmt.Sprintf("HTTP %d :: %s", res.StatusCode, string(data[:])))
	}

	return data, nil
}

func (b *BigIQ) iControlPath(parts []string) string {
	var buffer bytes.Buffer
	for i, p := range parts {
		buffer.WriteString(strings.Replace(p, "/", "~", -1))
		if i < len(parts)-1 {
			buffer.WriteString("/")
		}
	}
	return buffer.String()
}

//Generic delete
func (b *BigIQ) delete(path ...string) error {
	req := &APIRequest{
		Method: "delete",
		URL:    b.iControlPath(path),
	}

	_, callErr := b.APICall(req)
	return callErr
}

// checkError handles any errors we get from our API requests. It returns either the
// message of the error, if any, or nil.
func (b *BigIQ) checkError(resp []byte) error {
	if len(resp) == 0 {
		return nil
	}

	var reqError RequestError

	err := json.Unmarshal(resp, &reqError)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\n%s", err.Error(), string(resp[:])))
	}

	err = reqError.Error()
	if err != nil {
		return err
	}

	return nil
}

//Get a url and populate an entity. If the entity does not exist (404) then the
//passed entity will be untouched and false will be returned as the second parameter.
//You can use this to distinguish between a missing entity or an actual error.
func (b *BigIQ) getForEntity(e interface{}, path ...string) (error, bool) {
	req := &APIRequest{
		Method:      "get",
		URL:         b.iControlPath(path),
		ContentType: "application/json",
	}

	resp, err := b.APICall(req)
	if err != nil {
		var reqError RequestError
		json.Unmarshal(resp, &reqError)
		if reqError.Code == 404 {
			return nil, false
		}
		return err, false
	}

	err = json.Unmarshal(resp, e)
	if err != nil {
		return err, false
	}

	return nil, true
}

func (b *BigIQ) getForEntityNew(e interface{}, path ...string) (error, bool) {
	req := &APIRequest{
		Method:      "get",
		URL:         b.iControlPath(path),
		ContentType: "application/json",
	}

	resp, err := b.APICall(req)
	if err != nil {
		var reqError RequestError
		json.Unmarshal(resp, &reqError)
		return err, false
	}
	err = json.Unmarshal(resp, e)
	if err != nil {
		return err, false
	}
	return nil, true
}
