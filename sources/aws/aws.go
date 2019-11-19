package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/puppetlabs/mimosa/sources/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Query gathers intances data from AWS
func Query(config map[string]string) (map[string]common.MimosaData, error) {
	defer common.LogTiming(time.Now(), "aws.Query")

	// Validate config
	if config["region"] == "" {
		return nil, fmt.Errorf("Source configuration must specify a region")
	}
	if config["accessKey"] == "" {
		return nil, fmt.Errorf("Source configuration must specify an accessKey")
	}
	if config["secretKey"] == "" {
		return nil, fmt.Errorf("Source configuration must specify a secretKey")
	}

	// Query AWS for instances
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(config["region"]),
		Credentials: credentials.NewStaticCredentials(
			config["accessKey"],
			config["secretKey"],
			""),
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create an AWS session: %v", err)
	}
	svc := ec2.New(session)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Print(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Print(err.Error())
		}
		return nil, err
	}

	// Gather instances
	items := map[string]common.MimosaData{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			id := *instance.InstanceId
			state := *instance.State.Name
			if state == "terminated" {
				continue
			}
			data, err := json.Marshal(instance)
			if err != nil {
				return nil, err
			}
			items[id] = common.MimosaData{
				Version: "1.0",
				Typ:     "aws-instance",
				Data:    data,
			}
		}
	}
	return items, nil
}
