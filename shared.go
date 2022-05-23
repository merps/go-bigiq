package bigiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	uriShared       = "shared"
	uriLicensing    = "licensing"
	uriActivation   = "activation"
	uriRegistration = "registration"
	uriFileTransfer = "file-transfer"
	uriUploads      = "uploads"

	activationComplete   = "LICENSING_COMPLETE"
	activationInProgress = "LICENSING_ACTIVATION_IN_PROGRESS"
	activationFailed     = "LICENSING_FAILED"
	activationNeedEula   = "NEED_EULA_ACCEPT"
)

// Installs the given license.
func (b *BigIQ) InstallLicense(licenseText string) error {
	r := map[string]string{"licenseText": licenseText}
	return b.put(r, uriShared, uriLicensing, uriRegistration)
}

// Revoke license.
func (b *BigIQ) RevokeLicense() error {
	//r := map[string]string{"licenseText": licenseText}
	return b.delete(uriShared, uriLicensing, uriRegistration)
}

// Upload a file
func (b *BigIQ) UploadFile(f *os.File) (*Upload, error) {
	if strings.HasSuffix(f.Name(), ".iso") {
		err := fmt.Errorf("File must not have .iso extension")
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return b.Upload(f, info.Size(), uriShared, uriFileTransfer, uriUploads, info.Name())
}

// Upload a file from a byte slice
func (b *BigIQ) UploadBytes(data []byte, filename string) (*Upload, error) {
	r := bytes.NewReader(data)
	size := int64(len(data))
	return b.Upload(r, size, uriShared, uriFileTransfer, uriUploads, filename)
}

// Generic delete
func (b *BigIQ) deleteReq(path ...string) ([]byte, error) {
	req := &APIRequest{
		Method: "delete",
		URL:    b.iControlPath(path),
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

func (b *BigIQ) deleteReqBody(body interface{}, path ...string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "delete",
		URL:         b.iControlPath(path),
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// Upload a file read from a Reader
func (b *BigIQ) Upload(r io.Reader, size int64, path ...string) (*Upload, error) {
	client := &http.Client{
		Transport: b.Transport,
		Timeout:   b.ConfigOptions.APICallTimeout,
	}
	options := &APIRequest{
		Method:      "post",
		URL:         b.iControlPath(path),
		ContentType: "application/octet-stream",
	}
	var format string
	if strings.Contains(options.URL, "mgmt/") {
		format = "%s/%s"
	} else {
		format = "%s/mgmt/%s"
	}
	url := fmt.Sprintf(format, b.Host, options.URL)
	chunkSize := 512 * 1024
	var start, end int64
	for {
		// Read next chunk
		chunk := make([]byte, chunkSize)
		n, err := r.Read(chunk)
		if err != nil {
			return nil, err
		}
		end = start + int64(n)
		// Resize buffer size to number of bytes read
		if n < chunkSize {
			chunk = chunk[:n]
		}
		body := bytes.NewReader(chunk)
		req, _ := http.NewRequest(strings.ToUpper(options.Method), url, body)
		if b.Token != "" {
			req.Header.Set("X-F5-Auth-Token", b.Token)
		} else {
			req.SetBasicAuth(b.User, b.Password)
		}
		req.Header.Add("Content-Type", options.ContentType)
		req.Header.Add("Content-Range", fmt.Sprintf("%d-%d/%d", start, end-1, size))
		// Try to upload chunk
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		data, _ := ioutil.ReadAll(res.Body)
		if res.StatusCode >= 400 {
			if res.Header.Get("Content-Type") == "application/json" {
				return nil, b.checkError(data)
			}

			return nil, fmt.Errorf("HTTP %d :: %s", res.StatusCode, string(data[:]))
		}
		defer res.Body.Close()
		var upload Upload
		err = json.Unmarshal(data, &upload)
		if err != nil {
			return nil, err
		}
		start = end
		if start >= size {
			// Final chunk was uploaded
			return &upload, err
		}
	}
}

func (b *BigIQ) post(body interface{}, path ...string) error {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         b.iControlPath(path),
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

func (b *BigIQ) postReq(body interface{}, path ...string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         b.iControlPath(path),
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}
func (b *BigIQ) put(body interface{}, path ...string) error {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         b.iControlPath(path),
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

func (b *BigIQ) putReq(body interface{}, path ...string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         b.iControlPath(path),
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

func (b *BigIQ) patch(body interface{}, path ...string) error {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "patch",
		URL:         b.iControlPath(path),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

func (b *BigIQ) fastPatch(body interface{}, path ...string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "patch",
		URL:         b.iControlPath(path),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// jsonMarshal specifies an encoder with 'SetEscapeHTML' set to 'false' so that <, >, and & are not escaped. https://golang.org/pkg/encoding/json/#Marshal
// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// Helper to copy between transfer objects and model objects to hide the myriad of boolean representations
// in the iControlREST api. DTO fields can be tagged with bool:"yes|enabled|true" to set what true and false
// marshal to.
func marshal(to, from interface{}) error {
	toVal := reflect.ValueOf(to).Elem()
	fromVal := reflect.ValueOf(from).Elem()
	toType := toVal.Type()
	for i := 0; i < toVal.NumField(); i++ {
		toField := toVal.Field(i)
		toFieldType := toType.Field(i)
		fromField := fromVal.FieldByName(toFieldType.Name)
		if fromField.Interface() != nil && fromField.Kind() == toField.Kind() {
			toField.Set(fromField)
		} else if toField.Kind() == reflect.Bool && fromField.Kind() == reflect.String {
			switch fromField.Interface() {
			case "yes", "enabled", "true":
				toField.SetBool(true)
				break
			case "no", "disabled", "false", "":
				toField.SetBool(false)
				break
			default:
				return fmt.Errorf("Unknown boolean conversion for %s: %s", toFieldType.Name, fromField.Interface())
			}
		} else if fromField.Kind() == reflect.Bool && toField.Kind() == reflect.String {
			tag := toFieldType.Tag.Get("bool")
			switch tag {
			case "yes":
				toField.SetString(toBoolString(fromField.Interface().(bool), "yes", "no"))
				break
			case "enabled":
				toField.SetString(toBoolString(fromField.Interface().(bool), "enabled", "disabled"))
				break
			case "true":
				toField.SetString(toBoolString(fromField.Interface().(bool), "true", "false"))
				break
			}
		} else {
			return fmt.Errorf("Unknown type conversion %s -> %s", fromField.Kind(), toField.Kind())
		}
	}
	return nil
}

func toBoolString(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}
