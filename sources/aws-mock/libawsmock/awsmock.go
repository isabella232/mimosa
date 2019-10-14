package libawsmock

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dnaeon/go-vcr/cassette"
	yaml "gopkg.in/yaml.v2"
)

// Run the source
func Run(bucket string) error {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Query AWS for instances (pulled from a local file)
	instances, err := LoadMachinesFromCassette("aws.yaml")
	if err != nil {
		return err
	}
	log.Printf("Instances: %v\n", len(instances))

	// Write each instance to the bucket
	for _, instance := range instances {
		// Write the instance to
		object := *instance.InstanceId
		log.Printf("Object: %v\n", object)

		wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
		bs, err := json.Marshal(instance)
		if err != nil {
			return err
		}
		_, err = wc.Write(bs)
		if err != nil {
			return err
		}
		err = wc.Close()
		if err != nil {
			return err
		}
		log.Printf("Wrote: %s\n", *instance.InstanceId)
	}

	return nil
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
