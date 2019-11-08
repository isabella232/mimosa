package inventory

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type routerMessage struct {
	Bucket    string `json:"bucket"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Workspace string `json:"workspace"`
}

type host struct {
	Name      string `firestore:"name"`
	Hostname  string `firestore:"hostname"`
	IP        string `firestore:"ip"`
	State     string `firestore:"state"`
	Source    string `firestore:"source"`
	Timestamp string `firestore:"timestamp"`
}

type conversionFunc func([]byte) (*host, error)

type pubsubHandlerFunc func(ctx context.Context, m *pubsub.Message) error

func build(convert conversionFunc) pubsubHandlerFunc {

	// This is the pubsub handler
	return func(ctx context.Context, m *pubsub.Message) error {

		// Unmarshal the message
		var routerMessage routerMessage
		err := json.Unmarshal(m.Data, &routerMessage)
		if err != nil {
			return fmt.Errorf("failed to unmarshal router message: %v", err)
		}

		// FIXME Check version is supported
		if routerMessage.Version == "" {
			return fmt.Errorf("no version found in the router message: %v", err)
		}

		// Read object from the bucket
		client, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}
		rc, err := client.Bucket(routerMessage.Bucket).Object(routerMessage.Name).NewReader(ctx)
		if err != nil {
			return err
		}
		defer rc.Close()
		object, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}

		// Convert the object to a host
		host, err := convert(object)
		if err != nil {
			return err
		}
		if host.Name == "" {
			return fmt.Errorf("host must have a name: %v", host)
		}

		// Update host fields
		source := routerMessage.Bucket
		host.Source = source
		host.Timestamp = time.Now().Format(time.RFC3339)

		// Compute a deterministic hash to use as firestore ID
		sha := sha1.New()
		sha.Write([]byte(source))
		sha.Write([]byte(host.Name))
		id := hex.EncodeToString(sha.Sum(nil))

		// Write the doc to the "hosts" collection
		fc, err := firestore.NewClient(ctx, firestore.DetectProjectID)
		if err != nil {
			return err
		}
		_, err = fc.Collection("ws").Doc(routerMessage.Workspace).Collection("hosts").Doc(id).Set(context.Background(), host)
		if err != nil {
			return err
		}

		return err
	}

}
