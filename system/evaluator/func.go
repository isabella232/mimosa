package evaluator

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

type routerMessage struct {
	Bucket    string `json:"bucket"`
	Name      string `json:"name"`
	EventType string `json:"eventType"`
	Version   string `json:"version"`
	Workspace string `json:"workspace"`
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

	return nil

}
