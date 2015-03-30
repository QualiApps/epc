package epc

import (
	"bytes"
	"fmt"
	"encoding/xml"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"time"
	"crypto/hmac"
    "crypto/sha256"
    "crypto/tls"
    "encoding/base64"
    "reflect"
)

type (
	EPC struct {
		Auth      Auth
		ProjectID string
		Zone      string
		Image     string
		Shape     string
	}
	Auth struct {
		AccessUser, AccessToken string
	}
)

func GetAuth(accessUser, accessToken string) Auth {
	return Auth{accessUser, accessToken}
}

const (
	or2_api_url = "https://orchestration.epam.com/maestro2/api/cli"
  	api_version = "2.503.201"
  	cli_version = "2.503.201"

  	describeProjects = "describe-projects"
  	describeImages = "describe-images"
  	describeRegions = "describe-regions"
  	describeShapes = "describe-shapes"
  	describeInstance = "describe-instances"

  	runInstance = "run-instances"

  	statusOK = "200"
)

var (
	actionInstance = map[string]string{
		"start": "start-instances",
		"stop": "stop-instances",
		"reboot": "reboot-instances",
		"remove": "terminate-instances",
	}
)

func NewEPC(auth Auth, ProjectID, Zone, Image, Shape string) *EPC {
	return &EPC{
		Auth:      auth,
		ProjectID: ProjectID,
		Zone:      Zone,
		Image:     Image,
		Shape:     Shape,
	}
}

// Calls maestro api
func (e *EPC) maestroApiCall(v url.Values) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	client := &http.Client{Transport: transport}

	params := v.Encode()

	finalEndpoint := fmt.Sprintf("%s", or2_api_url)
	req, err := http.NewRequest("POST", finalEndpoint, bytes.NewBufferString(params))
	if err != nil {
		return &http.Response{}, fmt.Errorf("error creating request from client")
	}

	date := getDate()

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "epam-or2-go")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Maestro-api-version", api_version)
	req.Header.Add("Maestro-sdk-version", cli_version)
	req.Header.Add("Maestro-access-id", e.Auth.AccessUser)
	req.Header.Add("Maestro-date",  date)
	req.Header.Add("Maestro-authorization", ComputeHmac256(e.Auth.AccessUser, e.Auth.AccessToken, date, "POST", params))

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("client encountered error while doing the request: %s", err.Error())
		return resp, fmt.Errorf("client encountered error while doing the request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return resp, err
	}

	return resp, err
}

// Checks api, returns true or false
func (e *EPC) checkMaestroApi() (bool, error) {
	check := true
	values := url.Values{}
	_, err := e.maestroApiCall(values)
	if err != nil {
		 check = false
	}

	return check, err
}

// Create instance by params (include a key pairs name)
func (e *EPC) createInstance(key string) (Instance, error) {
	values := url.Values{}
	values.Set("action", runInstance)
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	values.Set("shape", e.Shape)
	values.Set("imageId", e.Image)
	values.Set("key-name", key)
	
	var Response CreateStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return Instance{}, err
	}

	return Response.Instances.Instance, err
}

// Retrieves available instance info by instance-id
func (e *EPC) getInstance(id string) (Instance, error) {
	var id_err error
	if id == "" {
		return Instance{}, id_err
	}

	values := url.Values{}
	values.Set("project", e.ProjectID)
	values.Set("instances", id)
	values.Set("action", describeInstance)
	
	var Response InstanceStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return Instance{}, err
	}

	return Response.Instances.Instance, err
}

// Performs some actions on instance (start, stop, reboot, remove)
func (e *EPC) instance(id string, action string, force bool) error {
	values := url.Values{}
	values.Set("action", actionInstance[action])
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	values.Set("instances", id)

	// for kill action
	if force {
		values.Set("force", "1")
	}

	_, err := e.maestroApiCall(values)

	if err != nil {
		return fmt.Errorf("Error trying API call to " + action + " instance: %s", err)
	}

	return err
}

// Creates key pair by machine name
func (e *EPC) makeKeyPair(key string, action string) (Key, error) {
	values := url.Values{}
	values.Set("action", action)
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	values.Set("key-name", key)

	var Response KeyStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return Key{}, err
	}

	return Response.Keys.Key, err
}

// Retrieves key
func (e *EPC) getKeyPair(key string, action string) (Key, error) {
	values := url.Values{}
	values.Set("action", action)
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	values.Set("key-name", key)

	var Response DescKeyStatus
	err := e.getContents(values, &Response, false)
	if err != nil {
		return Key{}, err
	}

	return Response.Keys.Key, err
}

// Deletes key pair by machine name
func (e *EPC) removeKeyPair(key string, action string) (Key, error) {
	values := url.Values{}
	values.Set("action", action)
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	values.Set("key-name", key)

	var Response KeyStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return Key{}, err
	}

	return Response.Keys.Key, err
}

// Retrieves available images
func (e *EPC) getImages() ([]Image, error) {
	values := url.Values{}
	values.Set("action", describeImages)
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	
	var Response Status
	err := e.getContents(values, &Response, true)
	if err != nil {
		return nil, err
	}

	return Response.Images.Image, err
}

// Retrieves available shapes
func (e *EPC) getShapes() ([]Shape, error) {
	values := url.Values{}
	values.Set("project", e.ProjectID)
	values.Set("region", e.Zone)
	values.Set("action", describeShapes)
	
	var Response ShapeStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return nil, err
	}

	return Response.Shapes.Shape, err
}

// Retrieves available projects
func (e *EPC) getProjects() ([]Project, error) {
	values := url.Values{}
	values.Set("action", describeProjects)
	
	var Response ProjectStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return nil, err
	}

	return Response.Projects.Project, err
}

// Retrieves available zones
func (e *EPC) getZones() ([]Region, error) {
	values := url.Values{}
	values.Set("project", e.ProjectID)
	values.Set("action", describeRegions)
	
	var Response ZoneStatus
	err := e.getContents(values, &Response, true)
	if err != nil {
		return nil, err
	}

	return Response.Regions.Region, err
}

// Makes an action by api
func (e *EPC) getContents(values url.Values, unmarshalledResponse interface{}, checkStatus bool) error {
	resp, err := e.maestroApiCall(values)
	if err != nil {
		return fmt.Errorf("Error trying API call - " + values.Get("action") + ": %s", err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("Error reading Epam Cloud response body")
	}

	if xml.Unmarshal(contents, unmarshalledResponse); err != nil {
		return fmt.Errorf("Error unmarshalling Epam Cloud response XML: %s", err)
	}

	if checkStatus {
		if reflect.ValueOf(unmarshalledResponse).Elem().FieldByName("Code").String() != statusOK {
			return fmt.Errorf(reflect.ValueOf(unmarshalledResponse).Elem().FieldByName("Message").String())
		}
	}

	return err
}

// Returns current time by format EEE, dd MMM yyyy HH:mm:ss z (GMT 0)
func getDate() string {
	now := time.Now().UTC()
	return strings.Replace(now.Format(time.RFC1123), "UTC", "GMT", 1)
}

// HMAC-SHA256 Algorithm
func ComputeHmac256(user_id string, user_token string, date string, http_method string, params string) string {
    key := []byte(make_key_string(user_token, date))
    h := hmac.New(sha256.New, key)
    h.Write([]byte(make_hash_string(http_method, user_id, date, params)))
    
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Make a key string by format token + date
func make_key_string(user_token string, date string) string {
	return user_token + date
}

// Make a hash
// http_method string - POST, GET, etc...
// user_id string - it's a user access token
// date string - current date
// params string - the POST params
func make_hash_string(http_method string, user_id string, date string, params string) string {
	s := []string{http_method, user_id, date, strings.Replace(params, "&", ":", -1)}
	return strings.Join(s, ":")
}