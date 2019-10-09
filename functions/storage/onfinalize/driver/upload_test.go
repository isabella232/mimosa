package driver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDocument(t *testing.T) {
	c, err := auth()
	require.NoError(t, err)

	createDocument(c, "markf-test-bucket", "h5.json", someJSON())
}

func someJSON() string {
	return `{"blah":"yada"}`
}
