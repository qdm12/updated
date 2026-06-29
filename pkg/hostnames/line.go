package hostnames

import (
	"strings"
)

type cleanLineFunc func(line string) (cleaned string, err error)

type checkLineFunc func(line string) (ok bool)

func preCleanLine(line string, customClean cleanLineFunc) (cleaned string, err error) {
	line = strings.ToLower(line)

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

	if customClean != nil {
		line, err = customClean(line)
	}

	return line, err
}

func isLineValid(line string, customIsLineValid checkLineFunc) (valid bool) {
	if line == "" {
		return false
	} else if customIsLineValid == nil {
		return true
	}
	return customIsLineValid(line)
}

func postCleanLine(line string, customClean cleanLineFunc) (cleaned string, err error) {
	line = strings.TrimSpace(line)
	if customClean != nil {
		line, err = customClean(line)
	}
	return line, err
}
