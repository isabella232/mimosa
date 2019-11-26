package reaper

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

type routerMessage struct {
	Bucket            string `json:"bucket"`
	Name              string `json:"name"`
	EventType         string `json:"event-type"`
	MimosaType        string `json:"mimosa-type"`
	MimosaTypeVersion string `json:"mimosa-type-version"`
	Workspace         string `json:"workspace"`
}

// Reap an instance
func Reap(ctx context.Context, m *pubsub.Message) error {

	// Unmarshall the router message
	var routerMessage routerMessage
	err := json.Unmarshal(m.Data, &routerMessage)
	if err != nil {
		return err
	}
	log.Printf("reaping: %v", routerMessage)

	return nil

}
