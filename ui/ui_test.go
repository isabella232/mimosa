package ui

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	hosts := []Host{}
	h := Host{Name: "AppServerOne", PublicDNS: "server.domain.com", PublicIP: "1.2.3.4"}
	hosts = append(hosts, h)
	buf := new(bytes.Buffer)
	err := merge(buf, hosts)
	require.NoError(t, err)
	require.Contains(t, buf.String(), "server.domain.com")
	require.Contains(t, buf.String(), "AppServerOne")
	require.Contains(t, buf.String(), "1.2.3.4")
}

func TestGetHosts(t *testing.T) {
	hosts, err := getHosts()
	require.NoError(t, err)
	require.Len(t, hosts, 2)
	fmt.Printf("hosts is %v", hosts)
}
