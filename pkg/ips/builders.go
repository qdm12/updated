package ips

import "context"

// BuildMalicious obtains lists of IP addresses from different web sources
// and returns a list of CIDR IP ranges of malicious IP addresses.
func (b *builder) BuildMalicious(ctx context.Context) (ips []string, err error) {
	sources := []sourceType{
		{
			url: "https://iplists.firehol.org/files/firehol_level1.netset",
			customIsLineValid: func(line string) bool {
				return line != "0.0.0.0/8"
			},
		},
		// sourceType{
		// 	url: "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt",
		// 	customPreCleanLine: func(line string) string {
		// 		found := b.verifier.SearchIPv4(line)
		// 		if len(found) == 0 {
		// 			return ""
		// 		}
		// 		return found[0]
		// 	},
		// },
	}
	return b.buildForSources(ctx, "malicious", sources)
}

// BuildIPsFromHostnames builds a list of IP addresses obtained by resolving
// some hostnames given.
func (b *builder) BuildIPsFromHostnames(hostnames []string) (ips []string) {
	b.logger.Info("finding IP addresses from %d hostnames", len(hostnames))
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
	for n := 0; n < len(hostnames); n++ {
		newIPs := <-ch
		ips = append(ips, newIPs...)
	}
	b.logger.Info("found %d IP addresses from %d hostnames", len(ips), len(hostnames))
	return ips
}
