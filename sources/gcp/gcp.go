package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/puppetlabs/mimosa/sources/gcp/common"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type sourceMessage struct {
	Data []byte `json:"data"`
}

// HandleMessage from the matching topic telling this source to run
func HandleMessage(ctx context.Context, m sourceMessage) error {
	return common.Collect(Query)
}

// Query gathers intances data from AWS
func Query(config map[string]string) (map[string][]byte, error) {
	defer common.LogTiming(time.Now(), "gcp.Query")

	// Validate the config
	//
	// FIXME
	// There should be a token in here to access the GCP project
	// For now assign your source service account the "compute viewer" permission here: https://console.cloud.google.com/iam-admin/iam
	//
	if config["project"] == "" {
		return nil, fmt.Errorf("Source configuration must specify a project")
	}
	if config["zone"] == "" {
		return nil, fmt.Errorf("Source configuration must specify a zone")
	}

	// Query GCP for VM instances
	computeClient, err := google.DefaultClient(context.Background(), compute.ComputeScope)
	if err != nil {
		return nil, fmt.Errorf("Cannot create the GCP client: %v", err)
	}
	computeService, err := compute.New(computeClient)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to GCP compute service: %v", err)
	}
	instances, err := computeService.Instances.List(config["project"], config["zone"]).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to list compute instances: %v", err)
	}

	items := map[string][]byte{}
	for _, instance := range instances.Items {
		id := fmt.Sprintf("%d", instance.Id)
		data, err := json.Marshal(instance)
		if err != nil {
			return nil, err
		}
		items[id] = data
	}

	return items, nil
}
