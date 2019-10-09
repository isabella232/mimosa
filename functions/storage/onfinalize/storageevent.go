// Package onfinalize handles the onfinalize event for a gcs bucket
package onfinalize

import (
	"context"
	"io/ioutil"

	// "encoding/json"
	"fmt"
	"log"

	// "net/http"
	"time"

	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/storage"
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

// HandleGCSEvent handles an GCSEvent and (hopefully, in the future) creates another document in another bucket
// gcloud functions deploy MFHelloGCS --runtime go111 --trigger-resource YOUR_TRIGGER_BUCKET_NAME --trigger-event google.storage.object.finalize
//
// gcloud functions deploy MFHelloGCS --runtime go111 --trigger-resource default --trigger-event google.storage.object.finalize
// ERROR: (gcloud.functions.deploy) OperationError: code=7, message=Insufficient permissions to (re)configure a trigger (permission denied for bucket default). Please, give owner permissions to the editor role of the bucket and try again.
//
// gcloud functions deploy MFHelloGCS --runtime go111 --trigger-resource not_a_real_bucket --trigger-event google.storage.object.finalize gets
//ERROR: (gcloud.functions.deploy) OperationError: code=3, message=Cloud Storage trigger bucket not_a_real_bucket not found
//
// this seems to work (though does not show up in console due to region mismatch?)
//
// does appear in logs i.e.
// see also:
// HandleGCSEvent should handle a finalize event and can be registered as follows
//
// (run in this folder)
// gcloud functions deploy HandleGCSEvent --runtime go111 --trigger-resource markf-test-bucket --trigger-event google.storage.object.finalize
//
// to check success, run
//
// gcloud beta functions list
// gcloud beta functions describe HandleGCSEvent
//
// logs can be read using
// gcloud beta functions logs read HandleGCSEvent
//
func HandleGCSEvent(ctx context.Context, e GCSEvent) error {
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
	log.Printf("b: %v\n", b)

	return nil
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
