package sources

import (
	"github.com/puppetlabs/mimosa/sources/aws"
	"github.com/puppetlabs/mimosa/sources/common"
	"github.com/puppetlabs/mimosa/sources/gcp"
	"github.com/puppetlabs/mimosa/sources/netscan"
	"github.com/puppetlabs/mimosa/sources/qualys"
)

// AWS source
var AWS = common.Build(aws.Query)

// GCP source
var GCP = common.Build(gcp.Query)

// NetScan source
var NetScan = netscan.NetScanPubSub

// NetScanIot to interperate the response from devices
var NetScanIot = common.Build_Iot()

// Qualys source
var Qualys = common.Build(qualys.Query)
