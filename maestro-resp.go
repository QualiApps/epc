package epc

import "encoding/xml"

// Describe shapes response
type CreateStatus struct {
	XMLName    xml.Name  `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Instances  Instances `xml:"run-instances-response"`
}

type InstanceStatus struct {
	XMLName    xml.Name  `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Instances  Instances `xml:"describe-instances-response"`
}

type Instances struct {
	Instance   Instance `xml:"instance"`
}

type Instance struct {
	Id    string `xml:"instanceID,attr"`
	Region    string `xml:"region,attr"`
	State    string `xml:"state,attr"`
	Cpu    string `xml:"cpu,attr"`
	Memory    string `xml:"memory,attr"`
	Requested    string `xml:"requested,attr"`
	Owner    string `xml:"owner,attr"`
	Image    string `xml:"image,attr"`
	Shape    string `xml:"shape,attr"`
	GuestOS  string `xml:"guestOS,attr"`
	PrivateIP  string `xml:"privateIP,attr"`
	InstanceName  string `xml:"instanceName,attr"`
}

// Describe images response
type Status struct {
	XMLName xml.Name    `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Images  Images `xml:"describe-images-response"`
}

type Images struct {
	Image   []Image `xml:"image"`
}

type Image struct {
	Id    string `xml:"id,attr"`
	Description    string `xml:"description,attr"`
	Group    string `xml:"group,attr"`
	State    string `xml:"state,attr"`
	Size_MB    string `xml:"size_MB,attr"`
}

// Describe zones response
type ZoneStatus struct {
	XMLName xml.Name    `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Regions  Regions `xml:"describe-regions-response"`
}

type Regions struct {
	Region   []Region `xml:"region"`
}

type Region struct {
	Id    string `xml:"id,attr"`
}

// Describe projects response
type ProjectStatus struct {
	XMLName xml.Name    `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Projects  Projects `xml:"describe-projects-response"`
}

type Projects struct {
	Project   []Project `xml:"project"`
}

type Project struct {
	Id    string `xml:"projectID,attr"`
}

// Describe shapes response
type ShapeStatus struct {
	XMLName xml.Name    `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Shapes  Shapes `xml:"describe-shapes-response"`
}

type Shapes struct {
	Shape   []Shape `xml:"shape"`
}

type Shape struct {
	Name    string `xml:"name,attr"`
	Cpu    string `xml:"cpu,attr"`
	MemoryMB    string `xml:"memory_MB,attr"`
}

// Create key response
type KeyStatus struct {
	XMLName xml.Name    `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Keys  Keys `xml:"create-keypair-response"`
}

type DescKeyStatus struct {
	XMLName xml.Name    `xml:"status"`
	Code       string    `xml:"code,attr"`
	Message    string    `xml:"message,attr"`
	Keys  Keys `xml:"describe-keypairs-response"`
}

type Keys struct {
	Key   Key `xml:"key"`
}

type Key struct {
	Name    string `xml:"name,attr"`
	Owner    string `xml:"owner,attr"`
	Project    string `xml:"project,attr"`
	PrivateKey    string `xml:"privateKey,attr"`
	PublicKey     string `xml:"publicKey,attr"`
}