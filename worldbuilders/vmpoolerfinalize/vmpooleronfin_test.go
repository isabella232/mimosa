package vmpoolerfinalize

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/johnmccabe/go-vmpooler/vm"

	"github.com/stretchr/testify/require"
)

func TestCanUnmarshal(t *testing.T) {
	var instance vm.VM
	b := []byte(someJSON())
	err := json.Unmarshal(b, &instance)
	require.NoError(t, err)
	log.Printf("instance: %s\n", b)
}

func TestMapInstance(t *testing.T) {
	var instance vm.VM
	b := []byte(someJSON())
	err := json.Unmarshal(b, &instance)
	require.NoError(t, err)
	actual := mapInstance(instance)
	require.NoError(t, err)
	require.Equal(t, "q7zwg4d1bmnedka", actual["name"])
	require.Equal(t, "10.16.115.130", actual["public_ip"])
	require.Equal(t, "q7zwg4d1bmnedka.delivery.puppetlabs.net", actual["public_dns"])
}
func someJSON() string {
	return `{"Hostname":"q7zwg4d1bmnedka","Fqdn":"q7zwg4d1bmnedka.delivery.puppetlabs.net","Ip":"10.16.115.130","State":"running","Running":0,"Lifetime":12,"Tags":{"created_by":"vmpooler_bitbar"},"Template":{"Id":"centos-8-x86_64-pixa3","Os":"centos","Osver":"8-x86_64","Arch":"pixa3"}}`
}
