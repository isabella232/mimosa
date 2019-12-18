// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"fmt"
	"log"

	b64 "encoding/base64"

	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))
	resp, err := sendCommand("mimosa-andy-260909", "us-central1", "edges", "netscan-1dbdb398-5b42-4556-87ad-3105965273ec", string(m.Data))
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)

	return nil
}

// sendCommand sends a command to a device listening for commands.
func sendCommand(projectID string, region string, registryID string, deviceID string, sendData string) (*cloudiot.SendCommandToDeviceResponse, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	client, err := cloudiot.New(httpClient)
	if err != nil {
		return nil, err
	}

	req := cloudiot.SendCommandToDeviceRequest{
		BinaryData: b64.StdEncoding.EncodeToString([]byte(sendData)),
	}

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)

	response, err := client.Projects.Locations.Registries.Devices.SendCommandToDevice(name, &req).Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}
