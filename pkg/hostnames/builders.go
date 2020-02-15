package hostnames

import (
	"net/http"
	"strings"
)

func BuildSurveillance(httpClient *http.Client) (hostnames []string, err error) {
	sources := []sourceType{
		sourceType{
			url: "https://raw.githubusercontent.com/dyne/domain-list/master/data/nsa",
		},
		sourceType{
			url: "https://raw.githubusercontent.com/Cauchon/NSABlocklist-pi-hole-edition/master/HOSTS%20(including%20excessive%20GOV%20URLs)",
		},
		sourceType{
			url: "https://raw.githubusercontent.com/CHEF-KOCH/NSABlocklist/master/HOSTS/HOSTS",
			customPreCleanLine: func(line string) string {
				return strings.TrimPrefix(line, "0.0.0.0 ")
			},
			customIsLineValid: func(line string) bool {
				return strings.HasPrefix(line, "127.0.0.1 ")
			},
		},
	}
	return buildForSources(httpClient, "surveillance", sources)
}

func BuildMalicious(httpClient *http.Client) (hostnames []string, err error) {
	sources := []sourceType{
		sourceType{
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
		sourceType{
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
	return buildForSources(httpClient, "malicious", sources)
}

func BuildAds(httpClient *http.Client) (hostnames []string, err error) {
	sources := []sourceType{
		sourceType{
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
		sourceType{
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
	return buildForSources(httpClient, "ads", sources)
}
