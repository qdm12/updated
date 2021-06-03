package ips

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type sourceType struct {
	url                 string
	customPreCleanLine  func(line string) string
	customIsLineValid   func(line string) bool
	customPostCleanLine func(line string) string
}

func (b *builder) buildForSources(ctx context.Context, title string, sources []sourceType) (ips []string, err error) {
	b.logger.Info("building %s IPs...", title)
	for _, source := range sources {
		newIPs, err := b.buildForSource(
			ctx, source.url,
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

var (
	ErrBadStatusCode = errors.New("bad HTTP status code")
)

func (b *builder) buildForSource(
	ctx context.Context, url string,
	customPreCleanLine func(line string) string,
	customIsLineValid func(line string) bool,
	customPostCleanLine func(line string) string,
) (ips []string, err error) {
	tStart := time.Now()
	b.logger.Debug("building IPs %s...", url)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := b.client.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("%w: %d %s", ErrBadStatusCode, response.StatusCode, response.Status)
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = preCleanLine(line, customPreCleanLine)
		if isLineValid(line, customIsLineValid) { //nolint:nestif
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
