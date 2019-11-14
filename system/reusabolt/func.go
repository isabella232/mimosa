package reusabolt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
)

type target struct {
	Workspace string `json:"workspace"`
	ID        string `json:"id"`
}

type payload struct {
	User        string `json:"user"`
	Hostname    string `json:"hostname"`
	KeyMaterial []byte `json:"keymaterial"`
}

// TriggerReusabolt runs Cloud Run functions in response to pubsub messages
func TriggerReusabolt(ctx context.Context, m *pubsub.Message) error {
	log.Printf("Received pubsub message: %s", m.Data)

	// Unmarshal the target
	var target target
	err := json.Unmarshal(m.Data, &target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %v", err)
	}

	// Bucket holding secrets
	berglasBucket := os.Getenv("MIMOSA_SECRETS_BUCKET")
	if len(berglasBucket) == 0 {
		log.Panic("MIMOSA_SECRETS_BUCKET cannot be empty")
	}
	log.Printf("Secrets bucket: %s", berglasBucket)

	// Try checking for a default key
	keyMaterial, err := berglas.Resolve(ctx, fmt.Sprintf("berglas://%s/default", berglasBucket))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found key material in berglas")
	if len(keyMaterial) == 0 {
		log.Panic("Key material cannot be empty")
	}
	if len(keyMaterial) < 100 {
		log.Print("Warning - Key material is suspiciously short")
	}

	// Use a specified project if there is one or detect if running inside GCP
	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		project = firestore.DetectProjectID
	}
	log.Printf("Project: %s", project)

	// Read the host data from firestore
	fc, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}
	host, err := fc.Collection("ws").Doc(target.Workspace).Collection("hosts").Doc(target.ID).Get(ctx)
	if err != nil {
		return err
	}

	// Construct the payload
	hostname, err := host.DataAt("hostname")
	if err != nil {
		return err
	}
	payload := payload{
		User:        "ubuntu",
		Hostname:    hostname.(string),
		KeyMaterial: keyMaterial,
	}

	// Marshal the payload
	bs, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Look up the service URL
	serviceURL := os.Getenv("MIMOSA_SERVICE_URL")
	if len(serviceURL) == 0 {
		log.Panic("Service URL cannot be empty")
	}
	log.Printf("Service URL: %s", serviceURL)

	// Run Cloud Run function
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		log.Fatalf("failed to get id token: %+v", err)
	}
	body := bytes.NewReader(bs)
	req, err := http.NewRequest("POST", serviceURL, body)
	if err != nil {
		log.Fatalf("failed to create POST request: %+v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to execute POST request: %+v", err)
	}
	defer response.Body.Close()

	// Check status
	if response.StatusCode != 200 {
		log.Fatalf("POST response did not return 200 status: %d", response.StatusCode)
	}

	// Read body
	bs, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("failed to read POST response body: %+v", err)
	}

	log.Printf("POST response body: %s", bs)

	// Unmarshal result
	var result map[string]interface{}
	err = json.Unmarshal(bs, &result)
	if err != nil {
		log.Panicf("Failed to unmarshal result: %v", err)
	}
	result["timestamp"] = time.Now().Format(time.RFC3339)

	// Write the doc to the "tasks" collection
	_, _, err = fc.Collection("ws").Doc(target.Workspace).Collection("tasks").Add(ctx, result)
	return err

}
