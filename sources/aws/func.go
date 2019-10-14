package aws

import (
	"context"
	"log"
	"os"

	"github.com/puppetlabs/mimosa/aws-source/libaws"
)

type sourceMessage struct {
	Data []byte `json:"data"`
}

// SourceSubscriber connects to AWS and gets instance data
func SourceSubscriber(ctx context.Context, m sourceMessage) error {
	bucket := os.Getenv("MIMOSA_GCP_BUCKET")
	if len(bucket) == 0 {
		log.Fatal("MIMOSA_GCP_BUCKET environment variable must be set")
	}
	return libaws.Run(bucket)
}
