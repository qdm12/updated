package ips

import (
	"net"
	"net/http"

	"github.com/qdm12/golibs/logging"

	"github.com/qdm12/golibs/verification"
)

// BuildMalicious obtains lists of IP addresses from different web sources
// and returns a list of CIDR IP ranges of malicious IP addresses
func BuildMalicious(httpClient *http.Client) (IPs []string, err error) {
	sources := []sourceType{
		sourceType{
			url: "https://iplists.firehol.org/files/firehol_level1.netset",
			customIsLineValid: func(line string) bool {
				return line != "0.0.0.0/8"
			},
		},
		sourceType{
			url: "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt",
			customPreCleanLine: func(line string) string {
				found := verification.SearchIPv4(line)
				if len(found) == 0 {
					return ""
				}
				return found[0]
			},
		},
	}
	return buildForSources(httpClient, "malicious", sources)
}

// BuildIPsFromHostnames builds a list of IP addresses obtained by resolving
// some hostnames given.
func BuildIPsFromHostnames(hostnames []string) (IPs []string) {
	logging.Infof("finding IP addresses from %d hostnames", len(hostnames))
	ch := make(chan []string)
	for _, hostname := range hostnames {
		// TODO with 100 workers
		go func(hostname string) {
			var IPs []string
			newIPs, err := net.LookupIP(hostname)
			if err != nil {
				logging.Debug(err.Error())
				ch <- nil
				return
			}
			for _, newIP := range newIPs {
				IPs = append(IPs, newIP.String())
			}
			ch <- IPs
		}(hostname)
	}
	N := len(hostnames)
	for N > 0 {
		select {
		case newIPs := <-ch:
			N--
			IPs = append(IPs, newIPs...)
		}
	}
	logging.Infof("found %d IP addresses from %d hostnames", len(IPs), len(hostnames))
	return IPs
}
