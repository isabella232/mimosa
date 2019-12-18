package netscan

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"encoding/base64"
	"encoding/json"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudiot/v1"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

func NetScanPubSub(ctx context.Context, m *pubsub.Message) error {
	projectID := os.Getenv("GCP_PROJECT")
	bucketName := string(m.Data)
	if !strings.HasPrefix(bucketName, "source-") {
		return fmt.Errorf("message must be a bucket name starting with 'source-': %v", bucketName)
	}

	// Create GCP client and get a handle on the bucket
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}
	bucket := client.Bucket(bucketName)

	// Load source config from the bucket
	var config map[string]string
	err = unmarshalFromBucket(bucket, "config.json", &config)
	if err != nil {
		return fmt.Errorf("Cannot read config.json: %v", err)
	}

	// Validate config
	if config["ipRange"] == "" {
		return fmt.Errorf("Source configuration must specify a url")
	}

	ipRange := config["ipRange"]

	log.Printf("Got message payload: %v and %v", string(ipRange), string(bucketName))

	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Failed to create http client: %v", err)
	}
	iotClient, err := cloudiot.New(httpClient)
	if err != nil {
		log.Fatalf("Failed to create iot client: %v", err)
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}
	req := cloudiot.SendCommandToDeviceRequest{
		BinaryData: base64.StdEncoding.EncodeToString(jsonBytes),
	}
	// Need a better way to handle definition of regioun and registry
	region := "us-central1"
	registry := "edges"

	deviceName := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, bucketName)

	response, err := iotClient.Projects.Locations.Registries.Devices.SendCommandToDevice(deviceName, &req).Do()
	if err != nil {
		log.Fatalf("failed to send command to device %v: %v", deviceName, err)
	}
	log.Printf("Response: %v", response)
	return nil
}

func unmarshalFromBucket(bucket *storage.BucketHandle, object string, v interface{}) error {
	rc, err := bucket.Object(object).NewReader(context.Background())
	if err != nil {
		return err
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
