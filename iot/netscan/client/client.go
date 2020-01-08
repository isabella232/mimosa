package client

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/tatsushid/go-fastping"
)

const (
	initialRunBatch    = 8 // number of concurrent pings in the initial run
	initialRunMaxRTT   = 100 * time.Millisecond
	finalRunsMaxRTT    = time.Second
	finalRunsMaxCycles = -1
	maxIPsPerRange     = 256 * 256
)

// API defines the interface that clients should
// implement to get resources
type API interface {
	ScanIPs() (<-chan *net.IPAddr, <-chan error)
}

// Client Implements the API
type Client struct {
	ipRange string
}

// NewClient cosntructs a new client using the parameters passed in
func NewClient(ipRange string) API {
	return Client{ipRange}
}

// ScanIPs pings IPs from Client.ipRanges and sends responsive nodes to a nodeIPs channel
// CAUTION - needs to be run as a root user - go-fastping implements ICMP ping using raw socket
func (client Client) ScanIPs() (<-chan *net.IPAddr, <-chan error) {

	nodeIPs, errorMessages := make(chan *net.IPAddr), make(chan error)

	go func() {
		defer close(nodeIPs)

		splitRegex, err := regexp.Compile("[,\\n\\s]+") // split on ",", " ", "\"
		if err != nil {
			errorMessages <- err
			return
		}

		items := trimSpacesInArrayItems(splitRegex.Split(client.ipRange, -1))

		var ipAddrs []*net.IPAddr
		for _, item := range items {
			if len(item) == 0 {
				continue
			}
			if strings.Contains(item, "/") {
				cidrIPs, err := IPsByCIDR(item)
				if err != nil {
					errorMessages <- err
				}
				for _, ip := range cidrIPs {
					ipAddrs = append(ipAddrs, ipToIPAddr(ip))
				}
			} else if strings.Contains(item, "-") {
				rangeIPs, err := IPsByRange(item)
				if err != nil {
					errorMessages <- err
				}
				for _, ip := range rangeIPs {
					ipAddrs = append(ipAddrs, ipToIPAddr(ip))
				}
			} else {
				ip := net.ParseIP(item)
				if ip == nil {
					errorMessages <- err
				} else {
					ipAddrs = append(ipAddrs, ipToIPAddr(ip))
				}
			}
		}

		notResolved := make(map[string]*net.IPAddr, len(ipAddrs))
		for _, ipAddr := range ipAddrs {
			notResolved[ipAddr.String()] = ipAddr
		}

		// Initial run - pings initialRunBatch number of IPs in parallel
		//             - creates node resources from the responding ones,
		//             - collects the ones that are not responding

		p := fastping.NewPinger()
		p.MaxRTT = initialRunMaxRTT
		onRecv, onIdle := make(chan *net.IPAddr), make(chan bool)
		p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
			onRecv <- addr
		}

		for i, cycles := 0, len(ipAddrs)/initialRunBatch+1; i < cycles; i++ {

			endInterval := (i + 1) * initialRunBatch
			if endInterval > len(ipAddrs) {
				endInterval = len(ipAddrs)
			}
			currentIPAddrs := ipAddrs[i*initialRunBatch : endInterval]

			for _, ipAddr := range currentIPAddrs {
				p.AddIPAddr(ipAddr)
			}

			p.OnIdle = func() {
				// init pinger
				for _, ipAddr := range currentIPAddrs {
					p.RemoveIPAddr(ipAddr)
				}
				onIdle <- true
			}

			p.RunLoop()

			resolved := 0

		loop:
			for {
				select {
				case ipAddr := <-onRecv:
					nodeIPs <- ipAddr
					p.RemoveIPAddr(ipAddr)
					delete(notResolved, ipAddr.String())
					resolved++
					if resolved >= initialRunBatch {
						break loop // all current IPs have been discovered
					}
				case <-onIdle:
					break loop
				case <-p.Done():
					if err := p.Err(); err != nil {
						errorMessages <- err
					}
					break loop
				}
			}
			p.Stop()
		}

		if len(notResolved) > 0 {

			// Final runs - pings all remaining unresolved IPs
			//            - if within MaxRTT timeframe some IPs respond the process is repeated
			//              until finalRunsMaxCycles are exhausted

			p = fastping.NewPinger()
			p.MaxRTT = finalRunsMaxRTT
			p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
				onRecv <- addr
			}
			p.OnIdle = func() {
				onIdle <- true
			}

			for _, ipAddr := range notResolved {
				p.AddIPAddr(ipAddr)
			}

			p.RunLoop()

			resolvedInCycle := 0
			cycles := 0

		floop:
			for {
				select {
				case ipAddr := <-onRecv:
					nodeIPs <- ipAddr
					p.RemoveIPAddr(ipAddr)
					delete(notResolved, ipAddr.String())
					resolvedInCycle++
					if len(notResolved) == 0 {
						break floop // all IPs have been discovered
					}
				case <-onIdle:
					cycles++
					if finalRunsMaxCycles > 0 && cycles >= finalRunsMaxCycles ||
						resolvedInCycle == 0 {
						break floop
					}
					resolvedInCycle = 0
				case <-p.Done():
					if err := p.Err(); err != nil {
						errorMessages <- err
					}
					break floop
				}
			}
			p.Stop()
		}
	}()

	return nodeIPs, errorMessages
}

func trimSpacesInArrayItems(items []string) []string {
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}

// IPsByCIDR returns list of IPs based on the given CIDR parameter
func IPsByCIDR(cidr string) ([]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	cnt := 0
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		cnt++
		if cnt > maxIPsPerRange {
			return ips[1:],
				fmt.Errorf(
					"CIDR %s exceeded maximal number of IPs (%d); IPs starting %s will be skipped",
					cidr, maxIPsPerRange, ip.String())
		}
		ips = append(ips, duplIP(ip))
	}

	if len(ips) == 1 {
		// xx.xx.xx.xx/32 - it specifically means one single host, and there is no 'network address' and 'broadcast address'
		return ips, nil
	}

	if len(ips) < 1 {
		return []net.IP{}, nil
	}

	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil

}

// IPsByRange returns list of IPs based on the given IP range parameter
func IPsByRange(ipRange string) ([]net.IP, error) {
	limits := trimSpacesInArrayItems(strings.Split(ipRange, "-"))
	if len(limits) != 2 {
		return nil, fmt.Errorf("Parse error - '%s' is not valid IP range", ipRange)
	}
	lowerLimit := net.ParseIP(limits[0])
	if lowerLimit == nil {
		return nil, invalidIPAddressError(limits[0])
	}
	upperLimit := net.ParseIP(limits[1])
	if upperLimit == nil {
		return nil, invalidIPAddressError(limits[1])
	}

	var ips []net.IP
	cnt := 0
	for ip := lowerLimit; bytes.Compare(upperLimit, ip) >= 0; inc(ip) {
		cnt++
		if cnt > maxIPsPerRange {
			return ips,
				fmt.Errorf(
					"IP range %s exceeded maximal number of IPs (%d); IPs starting %s will be skipped",
					ipRange, maxIPsPerRange, ip.String())
		}
		ips = append(ips, duplIP(ip))
	}
	return ips, nil
}

func duplIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func invalidIPAddressError(ipAddress string) error {
	return fmt.Errorf("Parse error - '%s' is not valid IP address", ipAddress)
}

func ipToIPAddr(ip net.IP) *net.IPAddr {
	return &net.IPAddr{IP: ip}
}
