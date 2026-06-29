package ips

import (
	"strings"
)

type cleanLineFunc func(line string) (cleaned string)

type checkLineFunc func(line string) (ok bool)

func preCleanLine(line string, preClean cleanLineFunc) string {
	// remove comment
	lineWithoutCommentBuilder := strings.Builder{}
	for _, r := range line {
		if r == '#' {
			break
		}
		lineWithoutCommentBuilder.WriteRune(r)
	}
	line = lineWithoutCommentBuilder.String()

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
