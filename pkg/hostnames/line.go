package hostnames

import (
	"strings"
)

type cleanLineFunc func(line string) (cleaned string)

type checkLineFunc func(line string) (ok bool)

func preCleanLine(line string, customClean cleanLineFunc) (cleaned string) {
	line = strings.ToLower(line)

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

	if customClean != nil {
		line = customClean(line)
	}

	return line
}

func isLineValid(line string, customIsLineValid checkLineFunc) (valid bool) {
	if line == "" {
		return false
	} else if customIsLineValid == nil {
		return true
	}
	return customIsLineValid(line)
}

func postCleanLine(line string, customClean cleanLineFunc) (cleaned string) {
	line = strings.TrimSpace(line)
	if customClean != nil {
		line = customClean(line)
	}
	return line
}
