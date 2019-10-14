package awsfinalize

import (
	"context"
	"io/ioutil"

	"encoding/json"
	"fmt"
	"log"
	"time"

	// "cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

// HandleInstance handles an GCSEvent and looks for an aws ec2 instance
// gcloud functions deploy HandleInstance --runtime go111 --trigger-resource markf-test-bucket --trigger-event google.storage.object.finalize
func HandleInstance(ctx context.Context, e GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Event ID: %v\n", meta.EventID)
	log.Printf("Event type: %v\n", meta.EventType)
	log.Printf("Bucket: %v\n", e.Bucket)
	log.Printf("File: %v\n", e.Name)
	log.Printf("Metageneration: %v\n", e.Metageneration)
	log.Printf("Created: %v\n", e.TimeCreated)
	log.Printf("Updated: %v\n", e.Updated)
	uri := fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name)
	log.Printf("URI: %v\n", uri)
	uri = fmt.Sprintf("https://storage.googleapis.com/%s/%s", e.Bucket, e.Name)
	log.Printf("URI: %v\n", uri)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}

	// can we do this?
	obj := client.Bucket(e.Bucket).Object(e.Name)
	log.Printf("obj: %v\n", obj)

	// or this?
	b, err := read(client, e.Bucket, e.Name)
	if err != nil {
		return err
	}
	log.Printf("b: %s\n", b)

	var instance ec2.Instance
	err = json.Unmarshal(b, &instance)
	if err != nil {
		return err
	}
	log.Printf("instance: %s\n", b)

	fc, err := firestore.NewClient(ctx, "lyra-proj")
	if err != nil {
		return err
	}
	i, err := mapInstance(instance)
	if err != nil {
		return err
	}
	hosts := fc.Collection("hosts")
	doc, result, err := hosts.Add(context.Background(), i)
	log.Printf("doc: %v\n", doc)
	log.Printf("result: %v\n", result)
	return err
}

func mapInstance(instance ec2.Instance) (map[string]interface{}, error) {
	m := map[string]interface{}{
		"name":  *instance.InstanceId,
		"since": *instance.LaunchTime,
	}
	setIfNotNull(m, "public_ip", instance.PublicIpAddress)
	setIfNotNull(m, "public_dns", instance.PublicDnsName)
	return m, nil
}

func setIfNotNull(m map[string]interface{}, key string, value *string) {
	if value == nil {
		return
	}
	m[key] = *value
}

// read is taken from here
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go
func read(client *storage.Client, bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// [START download_file]
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}
