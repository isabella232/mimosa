package inventory

import (
	"encoding/xml"
	"fmt"
)

type asset struct {
	ID                    string `xml:"ID"`
	IP                    string `xml:"IP"`
	TRACKINGMETHOD        string `xml:"TRACKING_METHOD"`
	OS                    string `xml:"OS"`
	DNS                   string `xml:"DNS"`
	NETBIOS               string `xml:"NETBIOS"`
	LASTSCANDATETIME      string `xml:"LAST_SCAN_DATETIME"`
	LASTVMSCANNEDDATE     string `xml:"LAST_VM_SCANNED_DATE"`
	LASTVMSCANNEDDURATION string `xml:"LAST_VM_SCANNED_DURATION"`
	DETECTIONLIST         struct {
		DETECTION []struct {
			QID                   string `xml:"QID"`
			TYPE                  string `xml:"TYPE"`
			SEVERITY              string `xml:"SEVERITY"`
			SSL                   string `xml:"SSL"`
			RESULTS               string `xml:"RESULTS"`
			STATUS                string `xml:"STATUS"`
			FIRSTFOUNDDATETIME    string `xml:"FIRST_FOUND_DATETIME"`
			LASTFOUNDDATETIME     string `xml:"LAST_FOUND_DATETIME"`
			TIMESFOUND            string `xml:"TIMES_FOUND"`
			LASTTESTDATETIME      string `xml:"LAST_TEST_DATETIME"`
			LASTUPDATEDATETIME    string `xml:"LAST_UPDATE_DATETIME"`
			ISIGNORED             string `xml:"IS_IGNORED"`
			ISDISABLED            string `xml:"IS_DISABLED"`
			LASTPROCESSEDDATETIME string `xml:"LAST_PROCESSED_DATETIME"`
			PORT                  string `xml:"PORT"`
			PROTOCOL              string `xml:"PROTOCOL"`
		} `xml:"DETECTION"`
	} `xml:"DETECTION_LIST"`
	LASTPCSCANNEDDATE string `xml:"LAST_PC_SCANNED_DATE"`
}

func convert(object []byte) (*affectedHost, map[string]bool, error) {
	var asset asset
	err := xml.Unmarshal(object, &asset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal asset: %v", err)
	}
	affectedHost := &affectedHost{
		ID:       asset.ID,
		Hostname: asset.IP,
		Name:     asset.DNS,
	}
	vulns := map[string]bool{}
	for _, vuln := range asset.DETECTIONLIST.DETECTION {
		vulns[vuln.QID] = true
	}
	return affectedHost, vulns, nil
}
