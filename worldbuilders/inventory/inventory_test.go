package inventory

import (
	"context"
	"encoding/json"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/require"
)

func TestHandleMessage(t *testing.T) {
	t.Skip("this is a test driver for dev rather than real unit test")

	/** Guide to running this test:
	Step 0: Close your IDE and set your GOOGLE_APPLICATION_CREDENTIALS
	Step 1: comment out t.Skip(...)
	Step 2: Create a bucket and add in your json file
	Step 3: Replace routerMessage's Bucket and Name with
			your new Bucket and Name
	Step 4: Check initial run uses EventType "OBJECT_FINALIZE"
	Step 5: Run test and check firestore (data should appear)
	Step 6: Replace EventType with "OBJECT_DELETE"
	Step 7: Run test and check firestore (data should be removed)
	*/
	routerMessage := routerMessage{
		Bucket:            "philips-tempbucket",
		Name:              "aws-example.json",
		EventType:         "OBJECT_FINALIZE",
		MimosaType:        "new",
		MimosaTypeVersion: "something",
		Workspace:         "philip",
	}

	handler := build(convertAWS)
	data, err := json.Marshal(routerMessage)
	require.NoError(t, err)
	err = handler(context.Background(), &pubsub.Message{Data: data})
	require.NoError(t, err)
}
