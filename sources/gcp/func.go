package gcp

import (
	"context"
	"log"
	"os"

	"github.com/puppetlabs/mimosa/sources/gcp/libgcp"
)

type sourceMessage struct {
	Data []byte `json:"data"`
}

// SourceSubscriber connects to GCP and gets instance data
func SourceSubscriber(ctx context.Context, m sourceMessage) error {
	bucket := os.Getenv("MIMOSA_GCP_BUCKET")
	if len(bucket) == 0 {
		log.Fatal("MIMOSA_GCP_BUCKET environment variable must be set")
	}
	return libgcp.Run(bucket)
}
