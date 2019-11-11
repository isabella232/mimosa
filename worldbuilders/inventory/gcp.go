package inventory

import (
	"encoding/json"
	"fmt"

	compute "google.golang.org/api/compute/v1"
)

// GCP  pubsub handler
var GCP = build(convertGCP)

func convertGCP(object []byte) (*host, error) {

	var instance compute.Instance
	err := json.Unmarshal(object, &instance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object")
	}

	// Check we have an ID since we're very dependent on it
	if instance.Id == 0 {
		return nil, fmt.Errorf("no compute instance ID could be found")
	}

	// Build host
	host := &host{
		Name:  fmt.Sprintf("%d", instance.Id),
		State: instance.Status,
	}
	if len(instance.NetworkInterfaces) > 0 && instance.NetworkInterfaces[0] != nil &&
		len(instance.NetworkInterfaces[0].AccessConfigs) > 0 && instance.NetworkInterfaces[0].AccessConfigs[0] != nil {
		host.IP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}

	return host, nil

}
