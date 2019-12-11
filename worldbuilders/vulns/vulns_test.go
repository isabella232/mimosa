package inventory

import (
	"context"
	"encoding/json"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/require"
)

func TestHandleMessage(t *testing.T) {

	t.Skip("this is a test driver rather than a real unit test")
	// set GOOGLE_APPLICATION_CREDENTIALS before running
	data, err := json.Marshal(&routerMessage{
		Bucket:            "scottyw",
		Name:              "qualys1.xml",
		EventType:         "OBJECT_FINALIZE",
		MimosaType:        "qualys-instance",
		MimosaTypeVersion: "1.0.0",
		Workspace:         "abcde",
	})
	require.NoError(t, err)
	err = HandleMessage(context.Background(), &pubsub.Message{Data: data})
	require.NoError(t, err)
	data, err = json.Marshal(&routerMessage{
		Bucket:            "scottyw",
		Name:              "qualys2.xml",
		EventType:         "OBJECT_FINALIZE",
		MimosaType:        "qualys-instance",
		MimosaTypeVersion: "1.0.0",
		Workspace:         "abcde",
	})
	require.NoError(t, err)
	err = HandleMessage(context.Background(), &pubsub.Message{Data: data})
	require.NoError(t, err)

	// Fail to make sure we see any output generated
	t.Fail()
}
