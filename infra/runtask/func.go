package runtask

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
)

// RunTask ...
func RunTask(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}

	auth := r.Header.Get("authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		w.Write([]byte(`No Firebase ID token was passed as a Bearer token in the Authorization header.
		Make sure you authorize your request by providing the following HTTP header:
		'Authorization: Bearer <Firebase ID Token>'`))
		w.WriteHeader(403)
		return
	}

	token, err := client.VerifyIDToken(ctx, auth[7:])
	if err != nil {
		log.Fatalf("error verifying ID token: %v", err)
	}
	log.Printf("Verified ID token: %v", token)

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
