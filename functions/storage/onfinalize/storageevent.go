// Package onfinalize handles the onfinalize event for a gcs bucket
package onfinalize

import (
	"context"
	"io/ioutil"

	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
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
	log.Printf("b: %s\n", b)

	var berry Berry
	err = json.Unmarshal(b, &berry)
	if err != nil {
		return err
	}
	log.Printf("berry: %s\n", b)
	log.Printf("berry.Flavors[0].Flavor.Name: %v\n", berry.Flavors[0].Flavor.Name)
	log.Printf("berry.Item.Name: %v\n", berry.Item.Name)

	//TODO can we get current projectID, do we want that
	fc, err := firestore.NewClient(ctx, "lyra-proj")
	hosts := fc.Collection("hosts")
	doc, result, err := hosts.Add(context.Background(), map[string]interface{}{
		"name": berry.Item.Name,
	})
	log.Printf("doc: %v\n", doc)
	log.Printf("result: %v\n", result)

	return err
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

// Berry came with the help of https://mholt.github.io/json-to-go/ and is based on data from
// curl https://pokeapi.co/api/v2/berry/1
type Berry struct {
	Firmness struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"firmness"`
	Flavors []struct {
		Flavor struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"flavor"`
		Potency int `json:"potency"`
	} `json:"flavors"`
	GrowthTime int `json:"growth_time"`
	ID         int `json:"id"`
	Item       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"item"`
	MaxHarvest       int    `json:"max_harvest"`
	Name             string `json:"name"`
	NaturalGiftPower int    `json:"natural_gift_power"`
	NaturalGiftType  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"natural_gift_type"`
	Size        int `json:"size"`
	Smoothness  int `json:"smoothness"`
	SoilDryness int `json:"soil_dryness"`
}
