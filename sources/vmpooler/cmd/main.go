package main

import (
	"log"
	"os"

	"github.com/puppetlabs/mimosa/sources/common"
	"github.com/puppetlabs/mimosa/sources/vmpooler"
)

func main() {

	//check that GOOGLE_APPLICATION_CREDENTIALS is set as this (json file) will be used to create a new storage.NewClient(ctx)
	//the project_id is configured in the json file, see: https://cloud.google.com/docs/authentication/getting-started
	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(value) == 0 {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable must be set")
	}

	err := common.Collect(vmpooler.Query)
	if err != nil {
		log.Fatalf("Source failed with error: %v", err)
	}

}
