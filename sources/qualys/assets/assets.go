package assets

import "encoding/xml"

type Vulnerable_Hosts_List struct {
	XMLName  xml.Name  `xml:"HOST_LIST_VM_DETECTION_OUTPUT,omitempty" json:"HOST_LIST_VM_DETECTION_OUTPUT,omitempty"`
	Response *Response `xml:"RESPONSE,omitempty" json:"RESPONSE,omitempty"`
}

type Response struct {
	XMLName   xml.Name   `xml:"RESPONSE,omitempty" json:"RESPONSE,omitempty"`
	Datetime  *Datetime  `xml:"DATETIME,omitempty" json:"DATETIME,omitempty"`
	Host_List *Host_List `xml:"HOST_LIST,omitempty" json:"HOST_LIST,omitempty"`
	Warning   *Warning   `xml:"WARNING,omitempty" json:"WARNING,omitempty"`
}

type Datetime struct {
	XMLName xml.Name `xml:"DATETIME,omitempty" json:"DATETIME,omitempty"`
	String  string   `xml:",chardata" json:",omitempty"`
}

type Host_List struct {
	XMLName xml.Name `xml:"HOST_LIST,omitempty" json:"HOST_LIST,omitempty"`
	Hosts   []*Host  `xml:"HOST,omitempty" json:"HOST,omitempty"`
}

type Warning struct {
	XMLName xml.Name `xml:"WARNING,omitempty" json:"WARNING,omitempty"`
	Code    *Code    `xml:"CODE,omitempty" json:"CODE,omitempty"`
	URL     *URL     `xml:"URL,omitempty" json:"URL,omitempty"`
}

type Host struct {
	XMLName  xml.Name `xml:"HOST,omitempty" json:"HOST,omitempty"`
	ID       *ID      `xml:"ID,omitempty" json:"ID,omitempty"`
	InnerXML string   `xml:",innerxml" json:",omitempty"`
}

type Code struct {
	XMLName xml.Name `xml:"CODE,omitempty" json:"CODE,omitempty"`
	Code    string   `xml:",chardata" json:",omitempty"`
}

type URL struct {
	XMLName xml.Name `xml:"URL,omitempty" json:"URL,omitempty"`
	URL     string   `xml:",chardata" json:",omitempty"`
}

type ID struct {
	XMLName xml.Name `xml:"ID,omitempty" json:"ID,omitempty"`
	ID      int32    `xml:",chardata" json:",omitempty"`
}
