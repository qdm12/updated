package hostnames

import (
	"context"
	"fmt"
	"strings"
)

// BuildMalicious builds a list of malicious hostnames from different sources.
func (b *Builder) BuildMalicious(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{
			url:      "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts",
			preClean: preCleanQuadZeroes,
		},
		// See https://github.com/blocklistproject/Lists
		{url: "https://blocklistproject.github.io/Lists/alt-version/abuse-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/fraud-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/malware-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/phishing-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/ransomware-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/scam-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/redirect-nl.txt"},
		{
			url:      "https://raw.githubusercontent.com/mitchellkrogza/Badd-Boyz-Hosts/master/hosts",
			preClean: preCleanQuadZeroes,
		},
		{
			url:      "https://raw.githubusercontent.com/bigdargon/hostsVN/master/option/hosts-VN",
			preClean: preCleanQuadZeroes,
		},
		{
			url: "https://urlhaus.abuse.ch/downloads/hostfile/",
			preClean: func(line string) (cleaned string, err error) {
				if line == "" {
					return "", nil
				}
				const expectedPrefix = "127.0.0.1	"
				if !strings.HasPrefix(line, expectedPrefix) {
					return "", fmt.Errorf("preClean: line prefix unexpected: %s", line)
				}
				return line[len(expectedPrefix):], nil
			},
		},
	}
	return b.buildForSources(ctx, "malicious", sources)
}

// BuildAds builds a list of ad hostnames from different sources.
func (b *Builder) BuildAds(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		// See https://github.com/blocklistproject/Lists
		{url: "https://blocklistproject.github.io/Lists/alt-version/ads-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/tracking-nl.txt"},
		{
			url:      "https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&mimetype=plaintext&useip=0.0.0.0",
			preClean: preCleanQuadZeroes,
		},
		{
			url:      "https://raw.githubusercontent.com/anudeepND/blacklist/refs/heads/master/adservers.txt",
			preClean: preCleanQuadZeroes,
		},
	}
	return b.buildForSources(ctx, "ads", sources)
}

func preCleanQuadZeroes(line string) (cleaned string, err error) {
	switch {
	case line == "":
		return "", nil
	case strings.HasPrefix(line, "0.0.0.0 "): // common case
		return line[8:], nil
	case strings.HasPrefix(line, "127.0.0.1 "), strings.HasPrefix(line, "::1 "): // ignore localhost
		return "", nil
	default:
		return "", fmt.Errorf("preCleanQuadZeroes: line prefix unexpected: %s", line)
	}
}
