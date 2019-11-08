package inventory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertVMPooler(t *testing.T) {
	bs := []byte(VMPoolerJSON)
	host, err := convertVMPooler(bs)
	require.NoError(t, err)
	require.Equal(t, "q7zwg4d1bmnedka", host.Name)
	require.Equal(t, "q7zwg4d1bmnedka.delivery.puppetlabs.net", host.Hostname)
	require.Equal(t, "10.16.115.130", host.IP)
	require.Equal(t, "running", host.State)
}

var VMPoolerJSON = `{
    "Hostname": "q7zwg4d1bmnedka",
    "Fqdn": "q7zwg4d1bmnedka.delivery.puppetlabs.net",
    "Ip": "10.16.115.130",
    "State": "running",
    "Running": 0,
    "Lifetime": 12,
    "Tags": {
        "created_by": "vmpooler_bitbar"
    },
    "Template": {
        "Id": "centos-8-x86_64-pixa3",
        "Os": "centos",
        "Osver": "8-x86_64",
        "Arch": "pixa3"
    }
}`
