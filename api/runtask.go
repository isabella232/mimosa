package runtask

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

type target struct {
	Workspace string `json:"workspace"`
	ID        string `json:"id"`
}

// RunTask API handler
func RunTask(w http.ResponseWriter, r *http.Request) {

	untrustedData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	var target target
	err = json.Unmarshal(untrustedData, &target)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to unmarshal body: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("target: %v", target)

	// Validate
	// FIXME we need to check this user has permission to run this task
	if target.Workspace == "" {
		http.Error(w, "workspace must be specified", http.StatusBadRequest)
		return
	}
	if target.ID == "" {
		http.Error(w, "id must be specified", http.StatusBadRequest)
		return
	}

	// We don't trust the data posted to us enough to repost the same bytes to Reusabolt
	// Let's remarshal the target which contains exactly the right fields and has been verified
	trustedData, err := json.Marshal(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal message: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the pubsub message to Reusabolt
	ctx := context.Background()
	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		http.Error(w, "GCP_PROJECT environment variable must be set", http.StatusInternalServerError)
		return
	}
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create pubsub client: %v", err), http.StatusInternalServerError)
		return
	}
	topic := client.TopicInProject("reusabolt", project)
	result := topic.Publish(ctx, &pubsub.Message{Data: trustedData})
	_, err = result.Get(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to publish a message to the 'reusabolt' topic: %v", err), http.StatusInternalServerError)
		return
	}

}
