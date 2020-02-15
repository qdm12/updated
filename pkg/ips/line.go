package ips

import (
	"strings"
)

func preCleanLine(line string, customPreCleanLine func(line string) string) string {
	// remove comment
	var lineWithoutComment string
	for _, r := range line {
		if r == '#' {
			break
		}
		lineWithoutComment += string(r)
	}
	line = lineWithoutComment

	line = strings.TrimSpace(line)

	if customPreCleanLine != nil {
		line = customPreCleanLine(line)
	}

	return line
}

func isLineValid(line string, customIsLineValid func(line string) bool) bool {
	if line == "" {
		return false
	} else if customIsLineValid != nil && !customIsLineValid(line) {
		return false
	}
	return true
}

func postCleanLine(line string, customPostCleanLine func(line string) string) string {
	line = strings.TrimSpace(line)
	if customPostCleanLine != nil {
		line = customPostCleanLine(line)
	}
	return line
}
