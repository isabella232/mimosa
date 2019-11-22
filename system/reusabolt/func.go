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

type task struct {
	Status    string `firestore:"status"`
	Timestamp string `firestore:"timestamp"`
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
		return err
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
	hostRef := fc.Collection("ws").Doc(target.Workspace).Collection("hosts").Doc(target.ID)
	host, err := hostRef.Get(ctx)
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

	// Write the task result placeholder
	taskRef, _, err := fc.Collection("ws").Doc(target.Workspace).Collection("tasks").Add(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	// Update the host with a task reference
	task := task{
		Status:    "running",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_, err = hostRef.Set(ctx, map[string]interface{}{
		"tasks": map[string]interface{}{
			taskRef.ID: task,
		},
	}, firestore.MergeAll)
	if err != nil {
		return err
	}

	// Run Cloud Run function
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return fmt.Errorf("failed to get id token: %+v", err)
	}
	body := bytes.NewReader(bs)
	req, err := http.NewRequest("POST", serviceURL, body)
	if err != nil {
		return fmt.Errorf("failed to create POST request: %+v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute POST request: %+v", err)
	}
	defer response.Body.Close()

	// Check status
	if response.StatusCode == 200 {
		task.Status = "success"
	} else {
		task.Status = "failure"
	}
	task.Timestamp = time.Now().Format(time.RFC3339)

	// Read body
	bs, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read POST response body: %+v", err)
	}
	log.Printf("POST response body: %s", bs)

	// Unmarshal result
	var result map[string]interface{}
	err = json.Unmarshal(bs, &result)
	if err != nil {
		return err
	}

	// Update the host with the result
	if result["error"] != nil {
		task.Status = "failure"
	}
	_, err = hostRef.Set(ctx, map[string]interface{}{
		"tasks": map[string]interface{}{
			taskRef.ID: task,
		},
	}, firestore.MergeAll)
	if err != nil {
		return err
	}

	// Write the results to the task doc
	_, err = taskRef.Set(ctx, result)
	return err

}
