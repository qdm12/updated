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
	preClean            func(line string) string
	checkLine           func(line string) bool
	customPostCleanLine func(line string) string
}

func (b *Builder) buildForSources(ctx context.Context, title string, sources []sourceType) (ips []string, err error) {
	b.logger.Infof("building %s IPs...", title)
	for _, source := range sources {
		newIPs, err := b.buildForSource(
			ctx, source.url,
			source.preClean,
			source.checkLine,
			source.customPostCleanLine,
		)
		if err != nil {
			return nil, fmt.Errorf("building from %s: %w", source.url, err)
		}
		ips = append(ips, newIPs...)
	}
	b.logger.Infof("built %s IPs: %d IP address lines fetched", title, len(ips))
	return ips, nil
}

var ErrBadStatusCode = errors.New("bad HTTP status code")

func (b *Builder) buildForSource(ctx context.Context, url string,
	preClean cleanLineFunc, checkLine checkLineFunc, postClean cleanLineFunc,
) (ips []string, err error) {
	b.logger.Debug("building IPs from " + url + "...")
	tStart := time.Now()

	content, err := getContent(ctx, b.client, url)
	if err != nil {
		return nil, fmt.Errorf("getting content: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	ips = make([]string, 0, len(lines))

	for _, line := range lines {
		line = preCleanLine(line, preClean)

		if !isLineValid(line, checkLine) {
			continue
		}

		line = postCleanLine(line, postClean)

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

		b.logger.Warn("Not an IP address nor an IP subnet: " + line)
	}

	b.logger.Info("built IPs from " + url + " during " + time.Since(tStart).String())

	return ips, nil
}

func getContent(ctx context.Context, client *http.Client, url string) (content []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("%w: %d %s", ErrBadStatusCode, response.StatusCode, response.Status)
	}

	content, err = io.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, err
	}

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	return content, nil
}

func netIPIsPrivate(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}
