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

//
// BUILD AND RUN THE CONTAINER
//
// docker build . -t gcr.io/PROJECT_ID/runner;docker run -a STDOUT -a STDERR -it --env PORT=8080 -p 8080:8080 gcr.io/PROJECT_ID/runner
//

//
// CURL TO TEST
//
// curl localhost:8080 --data-binary "@payload.json"
//

type payload struct {
	User        string `json:"user"`
	Hostname    string `json:"hostname"`
	KeyMaterial []byte `json:"keymaterial"`
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

	// Check payload
	if payload.User == "" {
		log.Panic("User cannot be empty")
	}
	if payload.Hostname == "" {
		log.Panic("Hostname cannot be empty")
	}
	if len(payload.KeyMaterial) == 0 {
		log.Panic("KeyMaterial cannot be empty")
	}
	if len(payload.KeyMaterial) < 100 {
		log.Print("Warning - KeyMaterial is suspiciously short")
	}

	// Debug
	log.Printf("User: %s", payload.User)
	log.Printf("Hostname: %s", payload.Hostname)

	// Write key file
	pemFile, err := ioutil.TempFile(".", "mimosa-key-")
	if err != nil {
		log.Panicf("Failed to create key file: %v", err)
	}
	_, err = pemFile.Write([]byte(payload.KeyMaterial))
	if err != nil {
		log.Panicf("Failed to write key file: %v", err)
	}

	// Run bolt
	cmd := exec.Command("bolt", "task", "run", "facts",
		"--format", "json",
		"--private-key", pemFile.Name(),
		"--no-host-key-check",
		"--user", payload.User,
		"--nodes", payload.Hostname)
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
