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

type storagePayload struct {
	Metadata struct {
		MimosaType        string `json:"mimosa-type"`
		MimosaTypeVersion string `json:"mimosa-type-version"`
	} `json:"metadata"`
}

type routerMessage struct {
	Bucket            string `json:"bucket"`
	Name              string `json:"name"`
	EventType         string `json:"event-type"`
	MimosaType        string `json:"mimosa-type"`
	MimosaTypeVersion string `json:"mimosa-type-version"`
	Workspace         string `json:"workspace"`
}

// Route responds to updates to Cloud Storage by enqueueing the changes in the right place for processing
func Route(ctx context.Context, m *pubsub.Message) error {

	// Unmarshall the storage payload
	var storagePayload storagePayload
	err := json.Unmarshal(m.Data, &storagePayload)
	if err != nil {
		return err
	}
	log.Printf("storage payload: %v", storagePayload)

	// Objects that do not have both a type and version are not routed
	if storagePayload.Metadata.MimosaType == "" || storagePayload.Metadata.MimosaTypeVersion == "" {
		return nil
	}

	// Get the workspace from the bucket label
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}
	attrs, err := storageClient.Bucket(m.Attributes["bucketId"]).Attrs(ctx)
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

	// Build the router message
	routerMessage := routerMessage{
		Bucket:            m.Attributes["bucketId"],
		Name:              m.Attributes["objectId"],
		EventType:         m.Attributes["eventType"],
		MimosaType:        storagePayload.Metadata.MimosaType,
		MimosaTypeVersion: storagePayload.Metadata.MimosaTypeVersion,
		Workspace:         workspace,
	}
	log.Printf("router message: %v", routerMessage)
	data, err := json.Marshal(routerMessage)
	if err != nil {
		return err
	}

	// Queue the message onto the topic associated with this type
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}
	topicName := "type-" + storagePayload.Metadata.MimosaType
	topic := client.TopicInProject(topicName, project)
	result := topic.Publish(ctx, &pubsub.Message{Data: data})
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	// Queue the message for evaluation also
	topic = client.TopicInProject("system-evaluator", project)
	result = topic.Publish(ctx, &pubsub.Message{Data: data})
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	return nil

}
