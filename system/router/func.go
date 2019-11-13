package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/storage/v1"
)

type routerMessage struct {
	Bucket    string `json:"bucket"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Workspace string `json:"workspace"`
}

// Route responds to updates to Cloud Storage by enqueueing the changes in the right place for processing
func Route(ctx context.Context, object *storage.Object) error {

	// GCP sets this environment variable for the Cloud Function
	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		return fmt.Errorf("GCP_PROJECT environment variable must be set")
	}

	// This one should be set at deployment time
	workspace := os.Getenv("MIMOSA_WORKSPACE")
	if len(workspace) == 0 {
		return fmt.Errorf("MIMOSA_WORKSPACE environment variable must be set")
	}

	// Check the type
	mimosaType := object.Metadata["mimosa-type"]
	if mimosaType == "" {
		log.Printf("Ignoring object with no mimosa-type: %s/%s", object.Bucket, object.Name)
		return nil
	}

	// Queue the message onto the topic associated with this type
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}
	topic := client.TopicInProject(mimosaType, project)
	routerMessage := routerMessage{
		Bucket:    object.Bucket,
		Name:      object.Name,
		Version:   object.Metadata["mimosa-type-version"],
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
