package ips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_preCleanLine(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		line               string
		customPreCleanLine func(line string) string
		cleanedLine        string
	}{
		"empty input": {"", func(line string) string { return line }, ""},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cleanedLine := preCleanLine(tc.line, tc.customPreCleanLine)
			assert.Equal(t, tc.cleanedLine, cleanedLine)
		})
	}
}
