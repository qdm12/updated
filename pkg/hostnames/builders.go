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
	}
	return b.buildForSources(ctx, "ads", sources)
}
