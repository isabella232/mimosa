package runtask

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunTask(t *testing.T) {

	t.Skip("this is a test driver rather than a real unit test")

	os.Setenv("GCP_PROJECT", "mimosa-256008")
	target := target{
		Workspace: "ce3zn",
		ID:        "tY_aVYrG42fNqbdvyWHY5bCL5TY",
		Name:      "package",
		Params: map[string]string{
			"name":   "openssl",
			"action": "upgrade",
		},
	}

	data, err := json.Marshal(target)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(data))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RunTask)
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "%v", rr)
	require.Equal(t, "", rr.Body.String(), "%v", rr)

	t.Fail()

}
