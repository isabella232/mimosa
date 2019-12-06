package inventory

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type routerMessage struct {
	Bucket            string `json:"bucket"`
	Name              string `json:"name"`
	EventType         string `json:"event-type"`
	MimosaType        string `json:"mimosa-type"`
	MimosaTypeVersion string `json:"mimosa-type-version"`
	Workspace         string `json:"workspace"`
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
		log.Printf("router message: %v", routerMessage)

		// Compute a deterministic hash to use as firestore ID
		id, err := generateDeterministicID(routerMessage.Bucket, routerMessage.Name)
		if err != nil {
			return err
		}

		project := os.Getenv("MIMOSA_GCP_PROJECT")
		if len(project) == 0 {
			project = firestore.DetectProjectID
		}
		fc, err := firestore.NewClient(ctx, project)
		if err != nil {
			return err
		}

		switch routerMessage.EventType {
		case "OBJECT_DELETE":
			err = delete(ctx, fc, id, routerMessage)
		case "OBJECT_FINALIZE":
			err = update(ctx, fc, convert, id, routerMessage)
		default:
			err = fmt.Errorf("unknown event type: %s", routerMessage.EventType)
		}
		return err
	}
}

func update(ctx context.Context, fc *firestore.Client, convert conversionFunc, id string, routerMessage routerMessage) error {
	// FIXME Check version is supported
	if routerMessage.MimosaTypeVersion == "" {
		return fmt.Errorf("no mimosa version found in the router message")
	}

	// Read object from the bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(routerMessage.Bucket).Object(routerMessage.Name)

	// NOTE: We need to tweak how routerMessage and metadata is checked
	// attrs, err := obj.Attrs(ctx)
	// if err != nil {
	// 	return err
	// }

	// NOTE: we could cross-check these values vs the router message (e.g. routerMessage.Version), or just use one mechanism
	// if _, ok := attrs.Metadata["mimosa-type"]; !ok {
	// 	return fmt.Errorf("missing metadata mimosa-type")
	// }
	// if _, ok := attrs.Metadata["mimosa-type-version"]; !ok {
	// 	return fmt.Errorf("missing metadata mimosa-type-version")
	// }

	rc, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	object, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	// // Convert the object to a host
	// // FIXME pass the mimosa-type and mimosa-type-version value to the convert func
	host, err := convert(object)
	if err != nil {
		return err
	}
	if host.Name == "" {
		return fmt.Errorf("host must have a name: %v", host)
	}
	host.Timestamp = time.Now().Format(time.RFC3339)
	// Write the doc to the "hosts" collection
	_, err = fc.Collection("ws").Doc(routerMessage.Workspace).Collection("hosts").Doc(id).Set(ctx, host)
	if err != nil {
		return err
	}
	return nil
}

func delete(ctx context.Context, fc *firestore.Client, id string, routerMessage routerMessage) error {
	// Write the doc to the "hosts" collection

	_, err := fc.Collection("ws").Doc(routerMessage.Workspace).Collection("hosts").Doc(id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func generateDeterministicID(bucketName, objectName string) (string, error) {
	// Compute a deterministic hash to use as firestore ID
	sha := sha1.New()
	_, err := sha.Write([]byte(bucketName))
	if err != nil {
		return "", err
	}
	_, err = sha.Write([]byte(objectName))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sha.Sum(nil)), nil
}
