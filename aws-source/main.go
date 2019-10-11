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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type sourceConfig struct {
	Region    string `json:"region,omitempty"`
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
}

func main() {

	// Pull GCP configuration from the environment
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable must be set")
	}
	bucket := os.Getenv("GCP_BUCKET")
	if bucket == "" {
		log.Fatal("GCP_BUCKET environment variable must be set")
	}

	// Create GCP client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Read source config from the predetermined config object in the bucket
	//
	// Add the config manually like this:
	// gsutil cp config.json gs://aws1-test-bucket
	//
	rc, err := client.Bucket(bucket).Object("config.json").NewReader(ctx)
	if err != nil {
		log.Fatalf("Cannot find the config object: %v", err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatalf("Cannot read the config object: %v", err)
	}
	var sourceConfig sourceConfig
	err = json.Unmarshal(data, &sourceConfig)
	if err != nil {
		log.Fatalf("Cannot unmarshal the config object: %v", err)
	}

	// Validate the config
	if sourceConfig.Region == "" {
		log.Fatal("Source configuration must specify a region")
	}
	if sourceConfig.AccessKey == "" {
		log.Fatal("Source configuration must specify an accessKey")
	}
	if sourceConfig.SecretKey == "" {
		log.Fatal("Source configuration must specify a secretKey")
	}

	// Query AWS for instances
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(sourceConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			sourceConfig.AccessKey,
			sourceConfig.SecretKey,
			""),
	})
	svc := ec2.New(session)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)

	// Write each instance to the bucket
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// Write the instance to
			object := fmt.Sprintf("%s", *instance.InstanceId)
			wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
			bs, err := json.Marshal(instance)
			if err != nil {
				log.Fatal(err)
			}
			_, err = wc.Write(bs)
			if err != nil {
				log.Fatal(err)
			}
			err = wc.Close()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Wrote: %s\n", *instance.InstanceId)
		}
	}

}
