package qualys

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/puppetlabs/mimosa/sources/common"
	"github.com/puppetlabs/mimosa/sources/qualys/assets"
)

// Query gathers intances data from Qualys
func Query(config map[string]string) (map[string]common.MimosaData, error) {
	defer common.LogTiming(time.Now(), "qualys.Query")

	// Validate config
	if config["url"] == "" {
		return nil, fmt.Errorf("Source configuration must specify a url")
	}
	if config["username"] == "" {
		return nil, fmt.Errorf("Source configuration must specify an username")
	}
	if config["password"] == "" {
		return nil, fmt.Errorf("Source configuration must specify a password")
	}

	// Query Qualys for instances and by qualys... read in blueray data from a bucket.. for now...
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// this we eventually be replaced by a call to Qualys itself
	rc, err := client.Bucket(fmt.Sprintf("%s-blueray", os.Getenv("GCP_PROJECT"))).Object("qualys.xml").NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	vulnsXML, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var vulnAssets assets.Vulnerable_Hosts_List
	err = xml.Unmarshal(vulnsXML, &vulnAssets)
	if err != nil {
		return nil, err
	}

	if vulnAssets.Response == nil {
		err = fmt.Errorf("Response field in Qualys response is empty.")
		return nil, err
	}

	if vulnAssets.Response.Host_List == nil {
		err = fmt.Errorf("Response.Host_List field in Qualys response is empty. Failed to update vulnerabilities")
		return nil, err
	}

	// Gather asset data
	items := map[string]common.MimosaData{}
	for _, host := range vulnAssets.Response.Host_List.Hosts {
		// we need the ID for the filename in the bucket
		id := fmt.Sprintf("%d", host.ID.ID)

		// we remove the ID field so it is ommitted. This is because
		// innerxml will already contain it, and we don't want duplicates
		host.ID = nil
		// kept as XML so we don't lose any info at this stage,
		// we'll normalise it later when adding it to the db
		data, err := xml.Marshal(host)
		if err != nil {
			return nil, err
		}

		items[id] = common.MimosaData{
			Version: "1.0",
			Typ:     "qualys-instance",
			Data:    data,
		}
	}

	return items, nil
}
