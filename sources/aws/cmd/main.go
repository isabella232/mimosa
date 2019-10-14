package main

import (
	"log"
	"os"

	"github.com/puppetlabs/mimosa/aws-source/libaws"
)

type sourceConfig struct {
	Region    string `json:"region,omitempty"`
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
}

func main() {

	//check that GOOGLE_APPLICATION_CREDENTIALS is set as this (json file) will be used to create a new storage.NewClient(ctx)
	//the project_id is configured in the json file, see: https://cloud.google.com/docs/authentication/getting-started
	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(value) == 0 {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable must be set")
	}

	bucket := os.Getenv("MIMOSA_GCP_BUCKET")
	if len(bucket) == 0 {
		log.Fatal("MIMOSA_GCP_BUCKET environment variable must be set")
	}

	err := libaws.Run(bucket)
	if err != nil {
		log.Fatalf("AWS source failed with error: %v", err)
	}

}
