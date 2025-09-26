package hostnames

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

var regexHostname = regexp.MustCompile(`([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9_])(\.([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9]))*`) //nolint:lll

func (b *Builder) buildForSources(ctx context.Context, title string,
	sources []sourceType,
) (hostnames []string, err error) {
	b.logger.Debugf("building %s hostnames...", title)
	uniqueHostnames := make(map[string]bool)
	totalHostnames := 0

	for _, source := range sources {
		newHostnames, err := b.buildForSource(ctx, source.url,
			source.preClean, source.checkLine, source.postClean)
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
		if !regexHostname.MatchString(hostname) {
			continue
		}
		sortedHostnames = append(sortedHostnames, hostname)
	}
	sortedHostnames.Sort()

	b.logger.Infof("built %s hostnames: %d fetched, %d unique", title, totalHostnames, sortedHostnames.Len())

	return sortedHostnames, nil
}

var ErrBadStatusCode = errors.New("bad HTTP status code")

func (b *Builder) buildForSource(ctx context.Context, url string,
	preClean cleanLineFunc, checkLine checkLineFunc, postClean cleanLineFunc,
) (hostnames []string, err error) {
	b.logger.Debug("building hostnames " + url + "...")
	tStart := time.Now()

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

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	hostnames = make([]string, 0, len(lines))
	for _, line := range lines {
		line = preCleanLine(line, preClean)
		if isLineValid(line, checkLine) {
			line = postCleanLine(line, postClean)
			hostnames = append(hostnames, line)
		}
	}

	b.logger.Infof("built hostnames %s during %s", url, time.Since(tStart))

	return hostnames, nil
}
