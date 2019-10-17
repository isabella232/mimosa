package libgcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type sourceConfig struct {
	Project string `json:"project,omitempty"`
	Zone    string `json:"zone,omitempty"`
}

// Run the source
func Run(bucket string) error {

	log.Printf("Accessing bucket %s ...", bucket)

	// Create GCP client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	rc, err := client.Bucket(bucket).Object("config.json").NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Cannot find the config object: %v", err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("Cannot read the config object: %v", err)
	}
	var sourceConfig sourceConfig
	err = json.Unmarshal(data, &sourceConfig)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal the config object: %v", err)
	}

	// Validate the config
	//
	// FIXME
	// There should be a token in here to access the GCP project
	// For now assign your source service account the "compute viewer" permission here: https://console.cloud.google.com/iam-admin/iam
	//
	if sourceConfig.Project == "" {
		return fmt.Errorf("Source configuration must specify a project")
	}
	if sourceConfig.Zone == "" {
		return fmt.Errorf("Source configuration must specify a zone")
	}

	// Query GCP for VM instances
	// computeClient, err := compute.NewClient(ctx)
	computeClient, err := google.DefaultClient(context.Background(), compute.ComputeScope)
	if err != nil {
		return fmt.Errorf("Cannot create the GCP client: %v", err)
	}

	computeService, err := compute.New(computeClient)
	if err != nil {
		return fmt.Errorf("Failed to connect to GCP compute service: %v", err)
	}

	instances, err := computeService.Instances.List(sourceConfig.Project, sourceConfig.Zone).Do()
	if err != nil {
		return fmt.Errorf("Failed to connect to list compute instances: %v", err)
	}

	// Write each instance to the bucket
	for _, instance := range instances.Items {
		// FIXME
		// Should be checking for changes via the "state.json" file - see libaws
		object := fmt.Sprintf("%d", instance.Id)
		wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
		bs, err := json.Marshal(instance)
		if err != nil {
			return err
		}
		_, err = wc.Write(bs)
		if err != nil {
			return err
		}
		err = wc.Close()
		if err != nil {
			return err
		}
		log.Printf("Wrote: %s\n", object)
	}

	log.Printf("Done")
	return nil
}
