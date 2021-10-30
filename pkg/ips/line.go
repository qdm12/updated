package ips

import (
	"strings"
)

type cleanLineFunc func(line string) (cleaned string)

type checkLineFunc func(line string) (ok bool)

func preCleanLine(line string, preClean cleanLineFunc) string {
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

	if preClean != nil {
		line = preClean(line)
	}

	return line
}

func isLineValid(line string, checkLine checkLineFunc) bool {
	if line == "" {
		return false
	} else if checkLine == nil {
		return true
	}
	return checkLine(line)
}

func postCleanLine(line string, postClean cleanLineFunc) string {
	line = strings.TrimSpace(line)
	if postClean != nil {
		line = postClean(line)
	}
	return line
}
