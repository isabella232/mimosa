package bindings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleMessage(t *testing.T) {
	t.Skip("this is a test driver rather than a real unit test")
	hostname := "ec2-52-209-165-226.eu-west-1.compute.amazonaws.com"
	te := TaskExecution{
		Name:    "facts",
		Targets: []string{"*"},
		Inventory: inventory{
			Nodes: []inventoryNode{
				inventoryNode{
					Name: hostname,
					Config: inventoryConfig{
						Transport: "ssh",
						SSH: &inventorySSH{
							User: "ubuntu",
							PrivateKey: &inventoryPrivateKey{
								KeyData: "-----BEGIN RSA PRIVATE KEY-----\r\nXXXXXXXX\r\n-----END RSA PRIVATE KEY-----",
							},
						},
					},
				},
			},
		},
	}
	r := NewReusabolt("http://localhost:4567")
	code, body, err := r.RunTask(te)
	require.Equal(t, 200, code)
	require.NotEmpty(t, body)
	require.NoError(t, err)
}
