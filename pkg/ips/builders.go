package ips

import (
	"context"
	"fmt"
	"strconv"
)

// BuildMalicious obtains lists of IP addresses from different web sources
// and returns a list of CIDR IP ranges of malicious IP addresses.
func (b *Builder) BuildMalicious(ctx context.Context) (ips []string, err error) {
	sources := []sourceType{
		{
			url: "https://iplists.firehol.org/files/firehol_level1.netset",
			checkLine: func(line string) (ok bool) {
				return line != "0.0.0.0/8"
			},
		},
		{
			url: "https://raw.githubusercontent.com/stamparm/ipsum/master/levels/2.txt",
			preClean: func(line string) string {
				found := b.verifier.SearchIPv4(line)
				if len(found) == 0 {
					return ""
				}
				return found[0]
			},
		},
	}

	return b.buildForSources(ctx, "malicious", sources)
}

// BuildIPsFromHostnames builds a list of IP addresses obtained by resolving
// some hostnames given.
func (b *Builder) BuildIPsFromHostnames(hostnames []string) (ips []string) {
	b.logger.Info("finding IP addresses from " +
		strconv.Itoa(len(hostnames)) + " hostnames...")

	ch := make(chan []string)
	for _, hostname := range hostnames {
		// TODO with 100 workers
		go func(hostname string) {
			var IPs []string
			newIPs, err := b.lookupIP(hostname)
			if err != nil {
				b.logger.Debug(err.Error())
				ch <- nil
				return
			}
			for _, newIP := range newIPs {
				IPs = append(IPs, newIP.String())
			}
			ch <- IPs
		}(hostname)
	}

	for range hostnames {
		newIPs := <-ch
		ips = append(ips, newIPs...)
	}

	b.logger.Info(fmt.Sprintf("found %d IP addresses from %d hostnames", len(ips), len(hostnames)))
	return ips
}
