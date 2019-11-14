package vmpooler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/johnmccabe/go-vmpooler/vm"
	"github.com/puppetlabs/mimosa/sources/common"
)

const vmpoolerEndpoint = "https://vmpooler.delivery.puppetlabs.net/api/v1"
const myToken = "fzssuuj7tjfmh7ppa7rty2tc16texrf6"

// HandleMessage from the matching topic telling this source to run
func HandleMessage(ctx context.Context, m *pubsub.Message) error {
	log.Printf("Received pubsub message: %s", m.ID)
	return common.Collect(Query)
}

// Query gathers intances data from vmpooler
func Query(config map[string]string) (map[string]common.MimosaData, error) {
	defer common.LogTiming(time.Now(), "vmpooler.Query")

	// Validate config
	// FIXME proper config support needed here
	// if config["token"] == "" {
	// 	return nil, fmt.Errorf("Source configuration must specify a region")
	// }

	// Query for vmpooler instances
	c := vm.NewClient(vmpoolerEndpoint, myToken)
	virtualmachines, err := c.GetAll()
	if err != nil {
		return nil, err
	}

	// Gather instances
	items := map[string]common.MimosaData{}
	for _, vm := range virtualmachines {
		// Zero out fields that change every time
		vm.Running = 0

		// Marshal
		id := vm.Hostname
		data, err := json.Marshal(vm)
		if err != nil {
			return nil, err
		}
		items[id] = common.MimosaData{
			Version: "1.0",
			Typ:     "aws-instance",
			Data:    data,
		}
	}

	return items, nil
}
