package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type payload struct {
	Name      string   `json:"name"`
	Targets   []string `json:"targets"`
	Inventory struct {
		Nodes []struct {
			Name   string `json:"name"`
			Config struct {
				Transport string `json:"transport"`
				SSH       struct {
					User       string `json:"user"`
					PrivateKey struct {
						KeyData []byte `json:"key-data"`
					} `json:"private-key"`
					HostKeyCheck bool `json:"host-key-check"`
				} `json:"ssh"`
			} `json:"config"`
		} `json:"nodes"`
	} `json:"inventory"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var must be specified")
	}
	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %s", addr)
	http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, req *http.Request) {

	// Unmarshal payload
	var payload payload
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Panicf("Failed to read payload: %v", err)
	}
	req.Body.Close()
	err = json.Unmarshal(bs, &payload)
	if err != nil {
		log.Panicf("Failed to unmarshal payload: %v", err)
	}

	user := payload.Inventory.Nodes[0].Config.SSH.User
	hostname := payload.Inventory.Nodes[0].Name
	keyMaterial := payload.Inventory.Nodes[0].Config.SSH.PrivateKey.KeyData

	// Check payload
	if user == "" {
		log.Panic("User cannot be empty")
	}
	if hostname == "" {
		log.Panic("Hostname cannot be empty")
	}
	if len(keyMaterial) == 0 {
		log.Panic("KeyMaterial cannot be empty")
	}
	if len(keyMaterial) < 100 {
		log.Print("Warning - KeyMaterial is suspiciously short")
	}

	// Debug
	log.Printf("User: %s", user)
	log.Printf("Hostname: %s", hostname)

	// Write key file
	pemFile, err := ioutil.TempFile(".", "mimosa-key-")
	if err != nil {
		log.Panicf("Failed to create key file: %v", err)
	}
	_, err = pemFile.Write(keyMaterial)
	if err != nil {
		log.Panicf("Failed to write key file: %v", err)
	}

	// Run bolt
	cmd := exec.Command("bolt", "task", "run", "facts",
		"--format", "json",
		"--private-key", pemFile.Name(),
		"--no-host-key-check",
		"--user", user,
		"--nodes", hostname)
	result, err := cmd.Output()
	if err != nil {
		log.Printf("bolt command exited with an error: %v", err)
		result, err = json.Marshal(map[string]interface{}{"error": err})
		if err != nil {
			log.Fatalf("Marshaling the error failed: %v", err)
		}
	}
	log.Printf("Result: %s", result)

	// Return a response
	_, err = w.Write(result)
	if err != nil {
		log.Panicf("Failed to write response: %v", err)
	}

}
