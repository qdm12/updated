package hostnames

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type sourceType struct {
	url                 string
	customPreCleanLine  func(line string) string
	customIsLineValid   func(line string) bool
	customPostCleanLine func(line string) string
}

func (b *builder) buildForSources(title string, sources []sourceType) (hostnames []string, err error) {
	b.logger.Debug("building %s hostnames...", title)
	uniqueHostnames := make(map[string]bool)
	var newHostnames []string
	var totalHostnames int
	for _, source := range sources {
		newHostnames, err = b.buildForSource(
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
		if !b.verifier.MatchHostname(hostname) {
			b.logger.Warn("hostname %q does not seem valid", hostname)
			continue
		}
		sortedHostnames = append(sortedHostnames, hostname)
	}
	sortedHostnames.Sort()
	b.logger.Info("built %s hostnames: %d fetched, %d unique", title, totalHostnames, sortedHostnames.Len())
	return sortedHostnames, nil
}

func (b *builder) buildForSource(
	URL string,
	customPreCleanLine func(line string) string,
	customIsLineValid func(line string) bool,
	customPostCleanLine func(line string) string,
) (hostnames []string, err error) {
	tStart := time.Now()
	b.logger.Debug("building hostnames %s...", URL)
	content, status, err := b.client.GetContent(URL)
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status for %q is %d", URL, status)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = preCleanLine(line, customPreCleanLine)
		if isLineValid(line, customIsLineValid) {
			line = postCleanLine(line, customPostCleanLine)
			hostnames = append(hostnames, line)
		}
	}
	b.logger.Info("built hostnames %s during %s", URL, time.Since(tStart))
	return hostnames, nil
}
