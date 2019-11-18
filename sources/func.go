package sources

import (
	"github.com/puppetlabs/mimosa/sources/aws"
	"github.com/puppetlabs/mimosa/sources/gcp"
	"github.com/puppetlabs/mimosa/sources/vmpooler"
)

// AWS source
var AWS = aws.HandleMessage

// GCP source
var GCP = gcp.HandleMessage

// VMPooler source
var VMPooler = vmpooler.HandleMessage
