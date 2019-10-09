package onfinalize

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"cloud.google.com/go/storage"
)

func TestRead(t *testing.T) {
	c, err := auth()
	require.NoError(t, err)
	b, err := read(c, "markf-test-bucket", "h1.json")
	require.NoError(t, err)
	s := string(b)
	require.Equal(t, "{\"blah\":\"yada\"}", s)

	//gsutil ls
	//gsutil ls gs://markf-test-bucket/
	//gsutil cat gs://markf-test-bucket/h1.json
}

func auth() (*storage.Client, error) {
	// ensure that GOOGLE_APPLICATION_CREDENTIALS is set correctly, see:
	// https://cloud.google.com/docs/authentication/production
	// for now, enforce presence
	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(value) == 0 {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set, cannot authenticate")
	}

	ctx := context.Background()
	return storage.NewClient(ctx)
}
