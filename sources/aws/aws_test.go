package aws

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
	err := common.Build(Query)(context.Background(), &pubsub.Message{Data: []byte("source-1fd41071-0569-4a53-9521-8a7e236a6ab1")})
	require.NoError(t, err)
}
