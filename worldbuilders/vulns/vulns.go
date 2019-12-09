package inventory

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type routerMessage struct {
	Bucket            string `json:"bucket"`
	Name              string `json:"name"`
	EventType         string `json:"event-type"`
	MimosaType        string `json:"mimosa-type"`
	MimosaTypeVersion string `json:"mimosa-type-version"`
	Workspace         string `json:"workspace"`
}

type vulnerability struct {
	ID    string                   `firestore:"id"`
	Name  string                   `firestore:"name"`
	Score string                   `firestore:"score"`
	Count int                      `firestore:"count"`
	Hosts map[string]*affectedHost `firestore:"hosts"`
}

type affectedHost struct {
	ID       string `firestore:"id"`
	Name     string `firestore:"name"`
	Hostname string `firestore:"hostname"`
}

// HandleMessage to find vulnerabilities
func HandleMessage(ctx context.Context, m *pubsub.Message) error {

	// Unmarshal the message
	var routerMessage routerMessage
	err := json.Unmarshal(m.Data, &routerMessage)
	if err != nil {
		return fmt.Errorf("failed to unmarshal router message: %v", err)
	}
	log.Printf("router message: %v", routerMessage)

	// FIXME Check version is supported
	if routerMessage.MimosaTypeVersion == "" {
		return fmt.Errorf("no version found in the router message: %v", err)
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

	// Get vulns from the object
	host, vulns, err := convert(object)
	if err != nil {
		return err
	}

	// Calculate the host ID deterministically
	hostID, err := generateDeterministicID(routerMessage.Bucket, routerMessage.Name)
	if err != nil {
		return err
	}

	// Firestore client
	project := os.Getenv("MIMOSA_GCP_PROJECT")
	if len(project) == 0 {
		project = firestore.DetectProjectID
	}
	fc, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}

	// Update each vuln to add this host
	for vulnID := range vulns {

		// FIXME this whole thing should be transactional

		// Find the vuln doc
		var doc *firestore.DocumentSnapshot
		var ref *firestore.DocumentRef
		var vulnerability vulnerability
		iter := fc.Collection("ws").Doc(routerMessage.Workspace).Collection("vulns").Where("id", "==", vulnID).Limit(1).Documents(ctx)
		doc, err = iter.Next()
		if err == iterator.Done {
			// Document doesn't exist
			ref = fc.Collection("ws").Doc(routerMessage.Workspace).Collection("vulns").NewDoc()
			vulnerability.ID = vulnID
			vulnerability.Name = "Vulnerability " + vulnID // FIXME pull from the vuln database by vulnID
			vulnerability.Score = "9.9"                    // FIXME pull from the vuln database by vulnID
		} else if err != nil {
			// This is a real error
			return err
		} else {
			ref = doc.Ref
			err = doc.DataTo(&vulnerability)
			if err != nil {
				return err
			}
		}

		// Add this host to the vuln and write back to Firestore if it is not already present
		if vulnerability.Hosts == nil {
			vulnerability.Hosts = map[string]*affectedHost{}
		}
		if vulnerability.Hosts[hostID] == nil {
			vulnerability.Hosts[hostID] = host
			_, err = ref.Set(ctx, &vulnerability)
			if err != nil {
				log.Printf("error: failed to updated vuln document %s: %v", vulnID, err)
			}
		}

	}

	return err
}

func generateDeterministicID(bucketName, objectName string) (string, error) {
	// Compute a deterministic hash to use as firestore ID
	sha := sha1.New()
	_, err := sha.Write([]byte(bucketName))
	if err != nil {
		return "", err
	}
	_, err = sha.Write([]byte(objectName))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}
