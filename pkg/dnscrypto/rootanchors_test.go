package dnscrypto

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseRootAnchors(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input []byte
		out   TrustAnchor
		err   error
	}{
		"root anchors 2019-10-04": {
			[]byte(""),
			TrustAnchor{},
			errors.New("EOF"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			out, err := parseRootAnchors(tc.input)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.out, out)
		})
	}
}
