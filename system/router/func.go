package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"

	"cloud.google.com/go/pubsub"
)

type storageEvent struct {
	Bucket   string `json:"bucket"`
	Name     string `json:"name"`
	Metadata struct {
		MimosaType        string `json:"mimosa-type"`
		MimosaTypeVersion string `json:"mimosa-type-version"`
	} `json:"metadata"`
	Workspace string `json:"workspace"`
}

type routerMessage struct {
	Bucket    string `json:"bucket"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Workspace string `json:"workspace"`
}

// Route responds to updates to Cloud Storage by enqueueing the changes in the right place for processing
func Route(ctx context.Context, m *pubsub.Message) error {

	// Unmarshall the storage event
	var storageEvent storageEvent
	err := json.Unmarshal(m.Data, &storageEvent)
	if err != nil {
		return err
	}
	log.Printf("storage message: %v", storageEvent)

	// Objects that do not have both a type and version are not routed
	if storageEvent.Metadata.MimosaType == "" || storageEvent.Metadata.MimosaTypeVersion == "" {
		return nil
	}

	// Get the workspace from the bucket label
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}
	attrs, err := storageClient.Bucket(storageEvent.Bucket).Attrs(ctx)
	if err != nil {
		return err
	}
	workspace, present := attrs.Labels["ws"]
	if !present {
		return fmt.Errorf("Bucket must have a 'ws' label indicating its workspace")
	}
	if workspace == "" {
		return fmt.Errorf("Bucket has a 'ws' label but the workspace is empty")
	}

	// GCP sets this environment variable for the Cloud Function
	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		return fmt.Errorf("GCP_PROJECT environment variable must be set")
	}

	// Queue the message onto the topic associated with this type
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}
	topicName := "type-" + storageEvent.Metadata.MimosaType
	topic := client.TopicInProject(topicName, project)
	routerMessage := routerMessage{
		Bucket:    storageEvent.Bucket,
		Name:      storageEvent.Name,
		Version:   storageEvent.Metadata.MimosaTypeVersion,
		Workspace: workspace,
	}
	log.Printf("router message: %v", routerMessage)
	data, err := json.Marshal(routerMessage)
	if err != nil {
		return err
	}
	result := topic.Publish(ctx, &pubsub.Message{Data: data})
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	return nil

}
