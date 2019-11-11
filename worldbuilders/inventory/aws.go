package inventory

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// AWS pubsub handler
var AWS = build(convertAWS)

func convertAWS(object []byte) (*host, error) {

	var instance ec2.Instance
	err := json.Unmarshal(object, &instance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object")
	}

	// Check we have an AWS ID since we're very dependent on it
	if instance.InstanceId == nil {
		return nil, fmt.Errorf("no AWS instance ID could be found")
	}

	// Build host
	host := &host{
		Name: *instance.InstanceId,
	}
	if instance.PublicDnsName != nil {
		host.Hostname = *instance.PublicDnsName
	}
	if instance.PublicIpAddress != nil {
		host.IP = *instance.PublicIpAddress
	}
	if instance.State != nil && instance.State.Name != nil {
		host.State = *instance.State.Name
	}

	return host, nil

}
