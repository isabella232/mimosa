package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleMessage(t *testing.T) {
	t.Skip("this is a test driver rather than a real unit test")
	//set GOOGLE_APPLICATION_CREDENTIALS and MIMOSA_GCP_BUCKET
	err := HandleMessage(context.Background(), sourceMessage{})
	require.NoError(t, err)
}
