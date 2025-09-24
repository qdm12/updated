// Package ips provides functionality to build and clean IP lists.
package ips

import (
	"net"
	"net/http"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

// Builder builds IP lists.
type Builder struct {
	client   *http.Client
	logger   logging.Logger
	verifier verification.Verifier
	lookupIP func(host string) ([]net.IP, error)
}

// New returns a new builder of IP lists.
func New(client *http.Client, logger logging.Logger) *Builder {
	return &Builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
		lookupIP: net.LookupIP,
	}
}
