package sources

import (
	"github.com/puppetlabs/mimosa/sources/aws"
	"github.com/puppetlabs/mimosa/sources/common"
	"github.com/puppetlabs/mimosa/sources/gcp"
	"github.com/puppetlabs/mimosa/sources/vmpooler"
	"github.com/puppetlabs/mimosa/sources/qualys"
)

// AWS source
var AWS = common.Build(aws.Query)

// GCP source
var GCP = common.Build(gcp.Query)

// VMPooler source
var VMPooler = common.Build(vmpooler.Query)

// Qualys source
var Qualys = common.Build(qualys.Query)
