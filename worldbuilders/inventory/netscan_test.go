package inventory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertNetScan(t *testing.T) {
	bs := []byte(NetScanJSON)
	host, err := convertNetScan(bs)
	require.NoError(t, err)
	require.Equal(t, "10.16.123.52", host.Name)
	require.Equal(t, "10.16.123.52", host.Hostname)
	require.Equal(t, "10.16.123.52", host.IP)
}

var NetScanJSON = `{
    "name":"10.16.123.52",
    "privateIPv4":"10.16.123.52",
    "privateIPv6":""
}`
