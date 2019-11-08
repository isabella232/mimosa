package inventory

import (
	"encoding/json"
	"fmt"

	"github.com/johnmccabe/go-vmpooler/vm"
)

// VMPooler pubsub handler
var VMPooler = build(convertVMPooler)

func convertVMPooler(object []byte) (*host, error) {

	var instance vm.VM
	err := json.Unmarshal(object, &instance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object")
	}

	// Check we have an VMPooler ID since we're very dependent on it
	if instance.Hostname == "" {
		return nil, fmt.Errorf("no VMPooler instance hostname could be found")
	}

	// Build host
	host := &host{
		Name:     instance.Hostname,
		Hostname: instance.Fqdn,
		IP:       instance.Ip,
		State:    instance.State,
	}

	return host, nil

}
