package hostnames

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
)

type sourceType struct {
	url                 string
	customPreCleanLine  func(line string) string
	customIsLineValid   func(line string) bool
	customPostCleanLine func(line string) string
}

func buildForSources(httpClient *http.Client, title string, sources []sourceType) (hostnames []string, err error) {
	logging.Debugf("building %s hostnames...", title)
	uniqueHostnames := make(map[string]bool)
	var newHostnames []string
	var totalHostnames int
	for _, source := range sources {
		newHostnames, err = buildForSource(
			httpClient,
			source.url,
			source.customPreCleanLine,
			source.customIsLineValid,
			source.customPostCleanLine,
		)
		if err != nil {
			return nil, err
		}
		for _, hostname := range newHostnames {
			totalHostnames++
			uniqueHostnames[hostname] = true
		}
	}
	var sortedHostnames sort.StringSlice
	for hostname := range uniqueHostnames {
		if !verification.MatchHostname(hostname) {
			logging.Warnf("hostname %q does not seem valid", hostname)
			continue
		}
		sortedHostnames = append(sortedHostnames, hostname)
	}
	sortedHostnames.Sort()
	logging.Infof("built %s hostnames: %d fetched, %d unique", title, totalHostnames, sortedHostnames.Len())
	return sortedHostnames, nil
}

func buildForSource(
	httpClient *http.Client,
	URL string,
	customPreCleanLine func(line string) string,
	customIsLineValid func(line string) bool,
	customPostCleanLine func(line string) string,
) (hostnames []string, err error) {
	tStart := time.Now()
	logging.Debugf("building hostnames %s...", URL)
	content, err := network.GetContent(httpClient, URL, network.GetContentParamsType{DisguisedUserAgent: true})
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = preCleanLine(line, customPreCleanLine)
		if isLineValid(line, customIsLineValid) {
			line = postCleanLine(line, customPostCleanLine)
			hostnames = append(hostnames, line)
		}
	}
	logging.Infof("built hostnames %s during %s", URL, time.Since(tStart))
	return hostnames, nil
}
