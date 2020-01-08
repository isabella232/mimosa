package client

import (
	"os"
	"testing"
)

var client API

func TestMain(m *testing.M) {

	os.Setenv("MOCK", "TRUE")

	// Create a new client
	client = NewClient("2607:f8b0:4003:c00::6a, 77.75.79.30 77.75.79.31\n77.75.79.32, 77.75.79.40/30, 77.75.79.50-77.75.79.55, 192.168.0.1, 192.168.0.2, 192.168.0.3")

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestScanIPs(t *testing.T) {

	nodeIPs, errors := client.ScanIPs()

	hostCnt := 0
	errCnt := 0

loop:
	for {
		select {
		case ip, more := <-nodeIPs:
			if !more {
				break loop
			}
			hostCnt++
			t.Logf("Host: %s\n", ip)
		case err := <-errors:
			errCnt++
			t.Logf("Error: %v\n", err)
		}
	}
	t.Log("Hosts: ", hostCnt)
	t.Log("Errors: ", errCnt)

	if hostCnt != 50 {
		t.Errorf("Expected %d nodeIPs but got %d", 50, hostCnt)
	}

	if errCnt != 0 {
		t.Errorf("Expected %d errors but got %d", 0, errCnt)
	}
}

func TestIPsByCIDR(t *testing.T) {
	cidrIPs, err := IPsByCIDR("77.75.79.40/30")
	if err != nil {
		t.Error(err.Error())
	} else if len(cidrIPs) != 2 ||
		cidrIPs[0].String() != "77.75.79.41" ||
		cidrIPs[1].String() != "77.75.79.42" {
		t.Errorf("CIDR 77.75.79.40/30 was resolved as %v but expected [77.75.79.41 77.75.79.42]", cidrIPs)
	}

	cidrIPs, err = IPsByCIDR("77.75.79.40/31")
	if err != nil {
		t.Error(err.Error())
	} else if len(cidrIPs) != 0 {
		t.Errorf("CIDR 77.75.79.40/31 was resolved as %v but expected []", cidrIPs)
	}

	cidrIPs, err = IPsByCIDR("77.75.79.40/32")
	if err != nil {
		t.Error(err.Error())
	} else if len(cidrIPs) != 1 {
		t.Errorf("CIDR 77.75.79.40/32 was resolved as %v but expected [77.75.79.40]", cidrIPs)
	}

	const expectedErrMsg = "CIDR 77.75.79.40/0 exceeded maximal number of IPs (65536); IPs starting 0.1.0.0 will be skipped"
	cidrIPs, err = IPsByCIDR("77.75.79.40/0")
	if err == nil {
		t.Errorf("CIDR 77.75.79.40/0 - missing warning for skipped IPs")
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("CIDR 77.75.79.40/0 - wrong error message: got \"%s\" but expected \"%s\"", err.Error(), expectedErrMsg)
		}
	}
	if len(cidrIPs) != (256*256 - 1) {
		t.Errorf("CIDR 77.75.79.40/0 - expecting 65535 IPs but got %d", len(cidrIPs))
	}
}

func TestIPsByRange(t *testing.T) {
	rangeIPs, err := IPsByRange("77.75.79.50-77.75.79.52")
	if err != nil {
		t.Error(err.Error())
	} else if len(rangeIPs) != 3 ||
		rangeIPs[0].String() != "77.75.79.50" ||
		rangeIPs[1].String() != "77.75.79.51" ||
		rangeIPs[2].String() != "77.75.79.52" {
		t.Errorf("IP range 77.75.79.50-77.75.79.52 was resolved as %v but expected [77.75.79.50 77.75.79.51 77.75.79.52]", rangeIPs)
	}

	rangeIPs, err = IPsByRange("77.75.79.50-77.75.79.49")
	if err != nil {
		t.Error(err.Error())
	} else if len(rangeIPs) != 0 {
		t.Errorf("IP range 77.75.79.50-77.75.79.49 was resolved as %v but expected []", rangeIPs)
	}

	const expectedErrMsg = "IP range 0.75.79.50-77.75.79.50 exceeded maximal number of IPs (65536); IPs starting 0.76.79.50 will be skipped"
	rangeIPs, err = IPsByRange("0.75.79.50-77.75.79.50")
	if err == nil {
		t.Errorf("IP range 0.75.79.50-77.75.79.50 - missing warning for skipped IPs")
	} else {
		if err.Error() != expectedErrMsg {
			t.Errorf("IP range 0.75.79.50-77.75.79.50 - wrong error message: got \"%s\" but expected \"%s\"", err.Error(), expectedErrMsg)
		}
	}
	if len(rangeIPs) != 256*256 {
		t.Errorf("IP range 0.75.79.50-77.75.79.50 - expecting 65536 IPs but got %d", len(rangeIPs))
	}
}
