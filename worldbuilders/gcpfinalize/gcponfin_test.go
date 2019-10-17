package gcpfinalize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	compute "google.golang.org/api/compute/v1"
)

func TestCanUnmarshal(t *testing.T) {
	var instance compute.Instance
	b := []byte(someJSON())
	err := json.Unmarshal(b, &instance)
	require.NoError(t, err)
	// log.Printf("instance: %s\n", b)
}

func TestMapInstance(t *testing.T) {
	var instance compute.Instance
	b := []byte(someJSON())
	err := json.Unmarshal(b, &instance)
	require.NoError(t, err)
	actual := mapInstance(instance)
	require.NoError(t, err)
	require.Equal(t, "6506532502499607010", actual["name"])
	require.Equal(t, "34.70.81.227", actual["public_ip"])
	require.Equal(t, nil, actual["public_dns"])
}

func someJSON() string {
	return `{
		"cpuPlatform": "Intel Haswell",
		"creationTimestamp": "2019-10-17T07:47:42.615-07:00",
		"disks": [
			{
				"autoDelete": true,
				"boot": true,
				"deviceName": "instance-1",
				"guestOsFeatures": [
					{
						"type": "VIRTIO_SCSI_MULTIQUEUE"
					}
				],
				"interface": "SCSI",
				"kind": "compute#attachedDisk",
				"licenses": [
					"https://www.googleapis.com/compute/v1/projects/debian-cloud/global/licenses/debian-9-stretch"
				],
				"mode": "READ_WRITE",
				"source": "https://www.googleapis.com/compute/v1/projects/mimosa-256008/zones/us-central1-a/disks/instance-1",
				"type": "PERSISTENT"
			}
		],
		"id": "6506532502499607010",
		"kind": "compute#instance",
		"labelFingerprint": "42WmSpB8rSM=",
		"machineType": "https://www.googleapis.com/compute/v1/projects/mimosa-256008/zones/us-central1-a/machineTypes/n1-standard-1",
		"metadata": {
			"fingerprint": "ugHUPM_JwJM=",
			"kind": "compute#metadata"
		},
		"name": "instance-1",
		"networkInterfaces": [
			{
				"accessConfigs": [
					{
						"kind": "compute#accessConfig",
						"name": "External NAT",
						"natIP": "34.70.81.227",
						"networkTier": "PREMIUM",
						"type": "ONE_TO_ONE_NAT"
					}
				],
				"fingerprint": "8qjEIxN3e9I=",
				"kind": "compute#networkInterface",
				"name": "nic0",
				"network": "https://www.googleapis.com/compute/v1/projects/mimosa-256008/global/networks/default",
				"networkIP": "10.128.0.2",
				"subnetwork": "https://www.googleapis.com/compute/v1/projects/mimosa-256008/regions/us-central1/subnetworks/default"
			}
		],
		"reservationAffinity": {
			"consumeReservationType": "ANY_RESERVATION"
		},
		"scheduling": {
			"automaticRestart": true,
			"onHostMaintenance": "MIGRATE"
		},
		"selfLink": "https://www.googleapis.com/compute/v1/projects/mimosa-256008/zones/us-central1-a/instances/instance-1",
		"serviceAccounts": [
			{
				"email": "126377560493-compute@developer.gserviceaccount.com",
				"scopes": [
					"https://www.googleapis.com/auth/devstorage.read_only",
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring.write",
					"https://www.googleapis.com/auth/servicecontrol",
					"https://www.googleapis.com/auth/service.management.readonly",
					"https://www.googleapis.com/auth/trace.append"
				]
			}
		],
		"status": "RUNNING",
		"tags": {
			"fingerprint": "42WmSpB8rSM="
		},
		"zone": "https://www.googleapis.com/compute/v1/projects/mimosa-256008/zones/us-central1-a"
	}`
}
