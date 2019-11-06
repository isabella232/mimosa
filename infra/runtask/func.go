package runtask

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

// RunTask ...
func RunTask(w http.ResponseWriter, r *http.Request) {

	firestoreID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("failed to read the POST body: %v", err)
	}

	if len(firestoreID) == 0 {
		log.Fatalf("firestore ID cannot be empty")
	}

	Publish(firestoreID)

}

// Publish a message to pubsub
func Publish(firestoreID []byte) {

	ctx := context.Background()

	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		log.Fatal("GCP_PROJECT environment variable must be set")
	}

	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("failed to create pubsub client: %v", err)
	}
	topic := client.TopicInProject("reusabolt", project)
	result := topic.Publish(ctx, &pubsub.Message{Data: firestoreID})
	_, err = result.Get(ctx)
	if err != nil {
		log.Fatalf("failed to publish a message to the 'reusabolt' topic: %v", err)
	}

}
