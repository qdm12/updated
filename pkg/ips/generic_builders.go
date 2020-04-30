package ips

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/golibs/network"
)

type sourceType struct {
	url                 string
	customPreCleanLine  func(line string) string
	customIsLineValid   func(line string) bool
	customPostCleanLine func(line string) string
}

func (b *builder) buildForSources(title string, sources []sourceType) (ips []string, err error) {
	b.logger.Info("building %s IPs...", title)
	for _, source := range sources {
		newIPs, err := b.buildForSource(
			source.url,
			source.customPreCleanLine,
			source.customIsLineValid,
			source.customPostCleanLine,
		)
		if err != nil {
			return nil, err
		}
		ips = append(ips, newIPs...)
	}
	b.logger.Info("built %s IPs: %d IP address lines fetched", title, len(ips))
	return ips, nil
}

func (b *builder) buildForSource(
	url string,
	customPreCleanLine func(line string) string,
	customIsLineValid func(line string) bool,
	customPostCleanLine func(line string) string,
) (ips []string, err error) {
	tStart := time.Now()
	b.logger.Debug("building IPs %s...", url)
	content, status, err := b.client.GetContent(url, network.UseRandomUserAgent())
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code for %q is %d", url, status)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = preCleanLine(line, customPreCleanLine)
		if isLineValid(line, customIsLineValid) {
			line = postCleanLine(line, customPostCleanLine)
			// check for single IP
			if IP := net.ParseIP(line); IP != nil {
				if !netIPIsPrivate(IP) {
					ips = append(ips, IP.String())
				}
				continue
			}
			// check for CIDR
			IP, CIDRPtr, err := net.ParseCIDR(line)
			if err == nil {
				if !netIPIsPrivate(IP) {
					ips = append(ips, CIDRPtr.String())
				}
				continue
			}
			b.logger.Warn("%q is not an IP address nor an IP subnet", line)
		}
	}
	b.logger.Info("built IPs %s during %s", url, time.Since(tStart))
	return ips, nil
}

func netIPIsPrivate(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	privateCIDRBlocks := [8]string{
		"127.0.0.0/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}
	for i := range privateCIDRBlocks {
		_, CIDR, _ := net.ParseCIDR(privateCIDRBlocks[i])
		if CIDR.Contains(ip) {
			return true
		}
	}
	return false
}
