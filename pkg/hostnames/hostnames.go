// Package hostnames provides functions to build lists of hostnames from various sources.
package hostnames

import (
	"net/http"

	"github.com/qdm12/golibs/verification"
)

// Builder builds hostnames lists.
type Builder struct {
	client   *http.Client
	logger   Logger
	verifier verification.Verifier
}

// Logger represents a minimal logger interface.
type Logger interface {
	Debug(s string)
	Debugf(format string, args ...any)
	Info(s string)
	Infof(format string, args ...any)
}

// New returns a new builder of hostnames lists.
func New(client *http.Client, logger Logger) *Builder {
	return &Builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
	}
}
