// Package hostnames provides functions to build lists of hostnames from various sources.
package hostnames

import (
	"net/http"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

// Builder builds hostnames lists.
type Builder struct {
	client   *http.Client
	logger   logging.Logger
	verifier verification.Verifier
}

// New returns a new builder of hostnames lists.
func New(client *http.Client, logger logging.Logger) *Builder {
	return &Builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
	}
}
