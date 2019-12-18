package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/puppetlabs/mimosa/system/evaluator/aws"
)

type routerMessage struct {
	Bucket            string `json:"bucket"`
	Name              string `json:"name"`
	EventType         string `json:"event-type"`
	MimosaType        string `json:"mimosa-type"`
	MimosaTypeVersion string `json:"mimosa-type-version"`
	Workspace         string `json:"workspace"`
}

var evaluaters = map[string]func(string) ([]string, error){
	"aws-instance": aws.EvaluateInstance,
	// "gcp-instance":      evalGCPInstance,
	//"netscan-instance": netscan.evalInstance,
}

// Evaluate rules against cloud resources
func Evaluate(ctx context.Context, m *pubsub.Message) error {

	// Unmarshall the router message
	var routerMessage routerMessage
	err := json.Unmarshal(m.Data, &routerMessage)
	if err != nil {
		return err
	}
	log.Printf("router message: %v", routerMessage)

	// Check whether we can evaluate
	evaluate, present := evaluaters[routerMessage.MimosaType]
	if !present {
		return nil
	}

	// Read object from the bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(routerMessage.Bucket).Object(routerMessage.Name)
	rc, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	object, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	// Evaluate where we need to send this message
	topics, err := evaluate(string(object))
	if err != nil {
		return err
	}
	for _, topic := range topics {
		log.Printf("publishing to action topic %s", topic)
		err := publish(ctx, m, topic)
		if err != nil {
			log.Printf("failed to publish router message to actions topic %s: %v", topic, err)
		}
	}

	return nil
}

func publish(ctx context.Context, m *pubsub.Message, topicName string) error {

	// GCP sets this environment variable for the Cloud Function
	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		return fmt.Errorf("GCP_PROJECT environment variable must be set")
	}

	// Queue the message onto the chosen topic
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return err
	}

	// Queue the message for evaluation also
	topic := client.TopicInProject(topicName, project)
	result := topic.Publish(ctx, m)
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	return nil

}
