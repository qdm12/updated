package ips

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/golibs/network"

	"github.com/qdm12/golibs/logging"
)

var privateCIDRs []*net.IPNet

func init() {
	privateCIDRBlocks := [8]string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}
	for _, privateCIDRBlock := range privateCIDRBlocks {
		_, CIDR, _ := net.ParseCIDR(privateCIDRBlock)
		privateCIDRs = append(privateCIDRs, CIDR)
	}
}

type sourceType struct {
	url                 string
	customPreCleanLine  func(line string) string
	customIsLineValid   func(line string) bool
	customPostCleanLine func(line string) string
}

func buildForSources(httpClient *http.Client, title string, sources []sourceType) (IPs []string, err error) {
	logging.Infof("building %s IPs...", title)
	for _, source := range sources {
		newIPs, err := buildForSource(
			httpClient,
			source.url,
			source.customPreCleanLine,
			source.customIsLineValid,
			source.customPostCleanLine,
		)
		if err != nil {
			return nil, err
		}
		IPs = append(IPs, newIPs...)
	}
	logging.Infof("built %s IPs: %d IP address lines fetched", title, len(IPs))
	return IPs, nil
}

func buildForSource(
	httpClient *http.Client,
	URL string,
	customPreCleanLine func(line string) string,
	customIsLineValid func(line string) bool,
	customPostCleanLine func(line string) string,
) (IPs []string, err error) {
	tStart := time.Now()
	logging.Debugf("building IPs %s...", URL)
	content, err := network.GetContent(httpClient, URL, network.GetContentParamsType{DisguisedUserAgent: true})
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = preCleanLine(line, customPreCleanLine)
		if isLineValid(line, customIsLineValid) {
			line = postCleanLine(line, customPostCleanLine)
			// check for single IP
			if IP := net.ParseIP(line); IP != nil {
				if !netIPIsPrivate(IP) {
					IPs = append(IPs, IP.String())
				}
				continue
			}
			// check for CIDR
			IP, CIDRPtr, err := net.ParseCIDR(line)
			if err == nil {
				if !netIPIsPrivate(IP) {
					IPs = append(IPs, CIDRPtr.String())
				}
				continue
			}
			logging.Warnf("%q is not an IP address nor an IP subnet", line)
		}
	}
	logging.Infof("built IPs %s during %s", URL, time.Since(tStart))
	return IPs, nil
}

func netIPIsPrivate(netIP net.IP) bool {
	for i := range privateCIDRs {
		if privateCIDRs[i].Contains(netIP) {
			return true
		}
	}
	return false
}
