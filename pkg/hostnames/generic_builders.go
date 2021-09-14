package hostnames

import (
	"context"
	"errors"
	"fmt"
	"io"
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

func (b *builder) buildForSources(ctx context.Context, title string,
	sources []sourceType) (hostnames []string, err error) {
	b.logger.Debug("building %s hostnames...", title)
	uniqueHostnames := make(map[string]bool)
	var newHostnames []string
	var totalHostnames int
	for _, source := range sources {
		newHostnames, err = b.buildForSource(
			ctx, source.url,
			source.customPreCleanLine,
			source.customIsLineValid,
			source.customPostCleanLine,
		)
		if err != nil {
			return nil, fmt.Errorf("building from %s: %w", source.url, err)
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

var (
	ErrBadStatusCode = errors.New("bad HTTP status code")
)

func (b *builder) buildForSource(
	ctx context.Context, url string,
	customPreCleanLine func(line string) string,
	customIsLineValid func(line string) bool,
	customPostCleanLine func(line string) string,
) (hostnames []string, err error) {
	tStart := time.Now()
	b.logger.Debug("building hostnames %s...", url)
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
		if isLineValid(line, customIsLineValid) {
			line = postCleanLine(line, customPostCleanLine)
			hostnames = append(hostnames, line)
		}
	}
	b.logger.Info("built hostnames %s during %s", url, time.Since(tStart))
	return hostnames, nil
}
