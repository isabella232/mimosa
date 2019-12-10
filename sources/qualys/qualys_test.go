package qualys

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"

	"github.com/puppetlabs/mimosa/sources/common"
	"github.com/stretchr/testify/require"
)

func TestHandleMessage(t *testing.T) {
	t.Skip("this is a test driver rather than a real unit test")
	// set GOOGLE_APPLICATION_CREDENTIALS before running
	err := common.Build(Query)(context.Background(), &pubsub.Message{Data: []byte("source-2602a9e2-2724-4bf8-8bb8-2b3248fd4c97")})
	require.NoError(t, err)
}
