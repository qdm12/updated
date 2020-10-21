package hostnames

import (
	"strings"
)

func (b *builder) BuildSurveillance() (hostnames []string, err error) {
	sources := []sourceType{
		{
			url: "https://raw.githubusercontent.com/dyne/domain-list/master/data/nsa",
		},
		{
			url: "https://raw.githubusercontent.com/Cauchon/NSABlocklist-pi-hole-edition/master/HOSTS%20(including%20excessive%20GOV%20URLs)", //nolint:lll
		},
	}
	return b.buildForSources("surveillance", sources)
}

func (b *builder) BuildMalicious() (hostnames []string, err error) {
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
				return !(strings.HasPrefix(line, "127.0.0.1 ") ||
					strings.HasPrefix(line, "255.255.255.255") ||
					strings.HasPrefix(line, "::1") ||
					strings.HasPrefix(line, "fe80::1") ||
					strings.HasPrefix(line, "ff00::0") ||
					strings.HasPrefix(line, "ff02::1") ||
					strings.HasPrefix(line, "ff02::2") ||
					strings.HasPrefix(line, "ff02::3"))
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
	return b.buildForSources("malicious", sources)
}

func (b *builder) BuildAds() (hostnames []string, err error) {
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
	return b.buildForSources("ads", sources)
}
