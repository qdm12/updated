package hostnames

import (
	"context"
	"strings"
)

func (b *Builder) BuildSurveillance(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{url: "https://raw.githubusercontent.com/dyne/domain-list/master/data/nsa"},
	}
	return b.buildForSources(ctx, "surveillance", sources)
}

func (b *Builder) BuildMalicious(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{
			url: "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
			preClean: func(line string) (cleaned string) {
				line = strings.TrimPrefix(line, "0.0.0.0 ")
				if line == "0.0.0.0" {
					line = "" // discard
				}
				return line
			},
			checkLine: func(line string) (ok bool) {
				return !hasAnyOfPrefixes(line, "127.0.0.1 ", "255.255.255.255",
					"::1", "fe80::1", "ff00::0", "ff02::1", "ff02::2", "ff02:")
			},
		},
		{
			url: "https://raw.githubusercontent.com/k0nsl/unbound-blocklist/master/blocks.conf",
			preClean: func(line string) string {
				line = strings.TrimPrefix(line, "local-zone: \"")
				line = strings.TrimSuffix(line, "\" redirect")
				return line
			},
			checkLine: func(line string) bool {
				return !strings.HasPrefix(line, "local-data: \"")
			},
		},
		// See https://github.com/blocklistproject/Lists
		{url: "https://blocklistproject.github.io/Lists/alt-version/abuse-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/fraud-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/tracking-nl.txt"},
	}
	return b.buildForSources(ctx, "malicious", sources)
}

func (b *Builder) BuildAds(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{
			url: "https://raw.githubusercontent.com/notracking/hosts-blocklists/master/domains.txt",
			checkLine: func(line string) bool {
				return !strings.HasSuffix(line, "/::")
			},
			postClean: func(line string) string {
				line = strings.TrimPrefix(line, "address=/")
				line = strings.TrimSuffix(line, "/0.0.0.0")
				return strings.TrimSuffix(line, ".")
			},
		},
		{
			url: "https://raw.githubusercontent.com/notracking/hosts-blocklists/master/hostnames.txt",
			checkLine: func(line string) bool {
				return !strings.HasPrefix(line, ":: ")
			},
			postClean: func(line string) string {
				line = strings.TrimPrefix(line, "0.0.0.0 ")
				return strings.TrimSuffix(line, ".")
			},
		},
		// See https://github.com/blocklistproject/Lists
		{url: "https://blocklistproject.github.io/Lists/alt-version/ads-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/malware-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/phishing-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/ransomware-nl.txt"},
		{url: "https://blocklistproject.github.io/Lists/alt-version/scam-nl.txt"},
	}
	return b.buildForSources(ctx, "ads", sources)
}

func hasAnyOfPrefixes(s string, prefixes ...string) (has bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
