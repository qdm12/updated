package hostnames

import (
	"context"
	"strings"
)

func (b *builder) BuildSurveillance(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{
			url: "https://raw.githubusercontent.com/dyne/domain-list/master/data/nsa",
		},
	}
	return b.buildForSources(ctx, "surveillance", sources)
}

func (b *builder) BuildMalicious(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{
			url: "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
			customPreCleanLine: func(line string) string {
				line = strings.TrimPrefix(line, "0.0.0.0 ")
				if line == "0.0.0.0" {
					line = ""
				}
				return line
			},
			customIsLineValid: func(line string) bool {
				return !strings.HasPrefix(line, "127.0.0.1 ") &&
					!strings.HasPrefix(line, "255.255.255.255") &&
					!strings.HasPrefix(line, "::1") &&
					!strings.HasPrefix(line, "fe80::1") &&
					!strings.HasPrefix(line, "ff00::0") &&
					!strings.HasPrefix(line, "ff02::1") &&
					!strings.HasPrefix(line, "ff02::2") &&
					!strings.HasPrefix(line, "ff02::3")
			},
		},
		{
			url: "https://raw.githubusercontent.com/k0nsl/unbound-blocklist/master/blocks.conf",
			customPreCleanLine: func(line string) string {
				line = strings.TrimPrefix(line, "local-zone: \"")
				line = strings.TrimSuffix(line, "\" redirect")
				return line
			},
			customIsLineValid: func(line string) bool {
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

func (b *builder) BuildAds(ctx context.Context) (hostnames []string, err error) {
	sources := []sourceType{
		{
			url: "https://raw.githubusercontent.com/notracking/hosts-blocklists/master/domains.txt",
			customIsLineValid: func(line string) bool {
				return !strings.HasSuffix(line, "/::")
			},
			customPostCleanLine: func(line string) string {
				line = strings.TrimPrefix(line, "address=/")
				line = strings.TrimSuffix(line, "/0.0.0.0")
				return strings.TrimSuffix(line, ".")
			},
		},
		{
			url: "https://raw.githubusercontent.com/notracking/hosts-blocklists/master/hostnames.txt",
			customIsLineValid: func(line string) bool {
				return !strings.HasPrefix(line, ":: ")
			},
			customPostCleanLine: func(line string) string {
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
