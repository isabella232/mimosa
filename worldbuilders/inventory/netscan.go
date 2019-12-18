package inventory

import (
	"encoding/json"
	"fmt"
)

type NetScanHost struct {
	Name        string `json:"name"`
	PrivateIPv4 string `json:"privateIPv4"`
	PrivateIPv6 string `json:"privateIPv6"`
}

// NetScan pubsub handler
var NetScan = build(convertNetScan)

func convertNetScan(object []byte) (*host, error) {
	var outputHost host
	var instance NetScanHost
	err := json.Unmarshal(object, &instance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object")
	}

	outputHost.Name = instance.Name

	// Check if it is IPv4 or IPv6
	if instance.PrivateIPv4 == "" {
		outputHost.Hostname = instance.PrivateIPv6
		outputHost.IP = instance.PrivateIPv6
	} else {
		outputHost.Hostname = instance.PrivateIPv4
		outputHost.IP = instance.PrivateIPv4
	}

	return &outputHost, nil

}
