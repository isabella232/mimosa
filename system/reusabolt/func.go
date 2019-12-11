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

// TriggerReusabolt runs Cloud Run functions in response to pubsub messages
func TriggerReusabolt(ctx context.Context, m *pubsub.Message) error {

	// Unmarshal the target
	var target target
	err := json.Unmarshal(m.Data, &target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %v", err)
	}
	log.Printf("target: %v", target)

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
	project := os.Getenv("MIMOSA_GCP_PROJECT")
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
	hostname, err := host.DataAt("hostname")
	if err != nil {
		return err
	}

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

	// Find the Reusabolt service and get a token
	serviceURL := os.Getenv("MIMOSA_SERVICE_URL")
	if len(serviceURL) == 0 {
		log.Panic("Service URL cannot be empty")
	}
	log.Printf("Service URL: %s", serviceURL)
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return fmt.Errorf("failed to get id token: %+v", err)
	}

	// Build the Reusabolt payload
	te := taskExecution{
		Name:    target.Name,
		Params:  target.Params,
		Targets: []string{"*"},
		Inventory: inventory{
			Nodes: []inventoryNode{
				inventoryNode{
					Name: hostname.(string),
					Config: inventoryConfig{
						Transport: "ssh",
						SSH: &inventorySSH{
							User: "ubuntu",
							PrivateKey: &inventoryPrivateKey{
								KeyData: string(keyMaterial),
							},
						},
					},
				},
			},
		},
	}
	data, err := json.Marshal(te)
	if err != nil {
		return fmt.Errorf("failed to marshal the Reusabolt payload: %+v", err)
	}

	// Run Cloud Run function
	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", serviceURL+"/v1/task", body)
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
	var result map[string]interface{}
	if response.StatusCode == 200 {

		task.Status = "success"

		// Read body
		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read POST response body: %+v", err)
		}
		log.Printf("Reusabolt response: %s", data)

		// Unmarshal the result
		err = json.Unmarshal(data, &result)
		if err != nil {
			return err
		}

	} else {
		task.Status = "failure"
	}
	task.Timestamp = time.Now().Format(time.RFC3339)

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

type target struct {
	Workspace string            `json:"workspace"`
	ID        string            `json:"id"`
	Name      string            `json:"name,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
}

type task struct {
	Status    string `firestore:"status"`
	Timestamp string `firestore:"timestamp"`
}

type taskExecution struct {
	Targets   []string          `json:"targets,omitempty"`
	Name      string            `json:"name,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
	Inventory inventory         `json:"inventory,omitempty"`
}

type inventory struct {
	Nodes []inventoryNode `json:"nodes,omitempty"`
}

type inventoryNode struct {
	Name   string          `json:"name,omitempty"`
	Config inventoryConfig `json:"config,omitempty"`
}

type inventoryConfig struct {
	Transport string          `json:"transport"`
	WinRM     *inventoryWinRM `json:"winrm,omitempty"`
	SSH       *inventorySSH   `json:"ssh,omitempty"`
}

type inventoryWinRM struct {
	User              string `json:"user,omitempty"`
	Password          string `json:"password,omitempty"`
	SSL               bool   `json:"ssl"`                 // Don't omit empty
	AllowHTTPFallback bool   `json:"allow_http_fallback"` // Attempt http if https fails
	SSLVerify         bool   `json:"ssl-verify"`          // Don't omit empty
}

type inventorySSH struct {
	User         string               `json:"user,omitempty"`
	Password     string               `json:"password,omitempty"`
	PrivateKey   *inventoryPrivateKey `json:"private-key,omitempty"`
	HostKeyCheck bool                 `json:"host-key-check"` // Don't omit empty
}

type inventoryPrivateKey struct {
	KeyData string `json:"key-data,omitempty"`
}
