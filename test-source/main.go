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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dnaeon/go-vcr/cassette"
	yaml "gopkg.in/yaml.v2"
)

//MIMOSA_GCP_BUCKET=your-bucket go run main.go
func main() {
	ctx := context.Background()

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

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Query AWS for instances (pulled from a local file)
	instances, err := LoadMachinesFromCassette("aws.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Instances: %v\n", len(instances))

	// Write each instance to the bucket
	for _, instance := range instances {
		// Write the instance to
		object := *instance.InstanceId
		fmt.Printf("Object: %v\n", object)

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

// LoadMachinesFromCassette marshalls the yaml and json data
// into a structure, so that it can be compared with what the AWS API returns.
func LoadMachinesFromCassette(cassetteFile string) ([]*ec2.Instance, error) {
	data, err := ioutil.ReadFile(cassetteFile)
	if err != nil {
		return nil, err
	}
	var cassette cassette.Cassette
	err = yaml.Unmarshal(data, &cassette)
	if err != nil {
		return nil, err
	}

	machines := []*ec2.Instance{}

	for _, interaction := range cassette.Interactions {
		if interaction.Request.Form.Get("Action") == "DescribeInstances" {
			// GOT ONE.....
			results := ec2.DescribeInstancesOutput{}
			decoder := xml.NewDecoder(strings.NewReader(interaction.Response.Body))
			err := xmlutil.UnmarshalXML(&results, decoder, "")
			if err != nil {
				return nil, err
			}

			for _, inst := range results.Reservations {
				machines = append(machines, inst.Instances...)
			}
		}
	}

	return machines, nil
}
