package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"

	"github.com/puppetlabs/mimosa/infra/runner"
)

func main() {

	//check that GOOGLE_APPLICATION_CREDENTIALS is set as this (json file) will be used to create a new storage.NewClient(ctx)
	//the project_id is configured in the json file, see: https://cloud.google.com/docs/authentication/getting-started
	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(value) == 0 {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable must be set")
	}

	err := runner.WrapReusabolt(context.Background(), &pubsub.Message{Data: []byte("a8e1e136ae5ea7c143a345e99aae843f22d6e5b1")})
	if err != nil {
		log.Fatalf("Source failed with error: %v", err)
	}

}
