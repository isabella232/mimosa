package libaws

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type sourceConfig struct {
	Region    string `json:"region,omitempty"`
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
}

// Run the source
func Run(bucket string) error {

	log.Printf("Accessing bucket %s ...", bucket)

	// Create GCP client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Read source config from the predetermined config object in the bucket
	// NOTE: create a user and its access key with limited privileges as follows:
	//   aws iam create-user --user-name your-service-account
	//   aws iam create-access-key --user-name your-service-account
	//   (copy the AccessKeyId and SecretAccessKey into the accessKey and secretKey of the config.json file)
	//   aws iam attach-user-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess --user-name your-service-account
	// Add the config manually like this:
	// gsutil cp config.json gs://aws1-test-bucket
	//
	rc, err := client.Bucket(bucket).Object("config.json").NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Cannot find the config object: %v", err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("Cannot read the config object: %v", err)
	}
	var sourceConfig sourceConfig
	err = json.Unmarshal(data, &sourceConfig)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal the config object: %v", err)
	}

	// Validate the config
	if sourceConfig.Region == "" {
		return fmt.Errorf("Source configuration must specify a region")
	}
	if sourceConfig.AccessKey == "" {
		return fmt.Errorf("Source configuration must specify an accessKey")
	}
	if sourceConfig.SecretKey == "" {
		return fmt.Errorf("Source configuration must specify a secretKey")
	}

	// Query AWS for instances
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(sourceConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			sourceConfig.AccessKey,
			sourceConfig.SecretKey,
			""),
	})
	if err != nil {
		return fmt.Errorf("Cannot create an AWS session: %v", err)
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
		return err
	}

	// Write each instance to the bucket
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			object := *instance.InstanceId
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
	}

	log.Printf("Done")
	return nil
}
