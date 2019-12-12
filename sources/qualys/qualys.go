package qualys

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/puppetlabs/mimosa/sources/common"
	"github.com/puppetlabs/mimosa/sources/qualys/assets"
)

const (
	defaultURLString       = "%s/api/2.0/fo/asset/host/vm/detection/?action=list&truncation_limit=%d"
	defaultTruncationLimit = 500
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

	var vulnAssets []*assets.Host
	if config["blueray"] != "true" {
		fmt.Println("qualys")
		qualysAssets, err := getVulnDataFromQualys(config["url"], config["username"], config["password"])
		if err != nil {
			return nil, err
		}
		vulnAssets = qualysAssets
	} else {
		fmt.Println("blueray")
		bluerayAssets, err := getVulnDataFromBlueray()
		if err != nil {
			return nil, err
		}
		vulnAssets = bluerayAssets
	}

	// Gather asset data
	items := map[string]common.MimosaData{}
	for _, host := range vulnAssets {
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

func getVulnDataFromQualys(qualysApiURL, username, password string) ([]*assets.Host, error) {
	queryURL := fmt.Sprintf(defaultURLString, qualysApiURL, defaultTruncationLimit)
	return getVulnsFromQualysURL(queryURL, username, password)
}

func getVulnsFromQualysURL(queryURL, username, password string) ([]*assets.Host, error) {
	fmt.Println("in getVulnsFromQualysURL")
	var vulnAssets []*assets.Host
	requestHeaders := map[string]interface{}{"X-Requested-With": "Puppet Security Blanket"}

	request, err := http.NewRequest("GET", queryURL, bytes.NewReader([]byte{}))
	if err != nil {
		return vulnAssets, err
	}
	request.SetBasicAuth(username, password)

	for key, value := range requestHeaders {
		request.Header.Add(key, fmt.Sprintf("%v", value))
	}

	client := http.Client{Timeout: time.Second * 60}
	resp, err := client.Do(request)
	if err != nil {
		return vulnAssets, err
	}

	if resp.StatusCode != 200 {
		return vulnAssets, fmt.Errorf(resp.Status)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return vulnAssets, err
	}

	vulnsXML, err := ioutil.ReadAll(buf)
	if err != nil {
		return vulnAssets, err
	}
	vulnAssetsStruct, err := xmlBytesToStruct(vulnsXML)

	if vulnAssetsStruct.Response == nil {
		err := fmt.Errorf("Response field in Qualys response is empty.")
		return vulnAssets, err
	}

	if vulnAssetsStruct.Response.Host_List == nil {
		err := fmt.Errorf("Response.Host_List field in Qualys response is empty. Failed to update vulnerabilities")
		return vulnAssets, err
	}

	vulnAssets = vulnAssetsStruct.Response.Host_List.Hosts
	if vulnAssetsStruct.Response.Warning != nil &&
		vulnAssetsStruct.Response.Warning.Code != nil &&
		vulnAssetsStruct.Response.Warning.Code.Code == "1980" &&
		vulnAssetsStruct.Response.Warning.URL != nil {

		// this isn't needed, but saves some memory, which we may need for bigger calls
		vulnsXML = []byte{}
		additionalAssets, err := getVulnsFromQualysURL(vulnAssetsStruct.Response.Warning.URL.URL, username, password)
		if err != nil {
			return vulnAssets, err
		}
		vulnAssets = append(vulnAssets, additionalAssets...)
	}

	return vulnAssets, nil
}

func getVulnDataFromBlueray() ([]*assets.Host, error) {
	var vulnAssets []*assets.Host
	// Query Qualys for instances and by qualys... read in blueray data from a bucket.. for now...
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return vulnAssets, err
	}

	// this we eventually be replaced by a call to Qualys itself
	rc, err := client.Bucket(fmt.Sprintf("%s-blueray", os.Getenv("GCP_PROJECT"))).Object("qualys.xml").NewReader(ctx)
	if err != nil {
		return vulnAssets, err
	}
	defer rc.Close()

	vulnsXML, err := ioutil.ReadAll(rc)
	if err != nil {
		return vulnAssets, err
	}
	vulnAssetsStruct, err := xmlBytesToStruct(vulnsXML)

	if vulnAssetsStruct.Response == nil {
		err := fmt.Errorf("Response field in Qualys response is empty.")
		return vulnAssets, err
	}

	if vulnAssetsStruct.Response.Host_List == nil {
		err := fmt.Errorf("Response.Host_List field in Qualys response is empty. Failed to update vulnerabilities")
		return vulnAssets, err
	}

	return vulnAssetsStruct.Response.Host_List.Hosts, err
}

func xmlBytesToStruct(xmlBytes []byte) (assets.Vulnerable_Hosts_List, error) {
	var vulnAssets assets.Vulnerable_Hosts_List
	err := xml.Unmarshal(xmlBytes, &vulnAssets)
	return vulnAssets, err
}
