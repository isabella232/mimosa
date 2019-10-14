// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample buckets creates a bucket, lists buckets and deletes a bucket
// using the Google Storage API. More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.
package main

import (
	"log"
	"os"

	"github.com/puppetlabs/mimosa/test-source/libawsmock"
)

//MIMOSA_GCP_BUCKET=your-bucket go run main.go
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

	err := libawsmock.Run(bucket)
	if err != nil {
		log.Fatalf("AWS source failed with error: %v", err)
	}

}
