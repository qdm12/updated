// Package ips provides functionality to build and clean IP lists.
package ips

import (
	"net"
	"net/http"
)

// Builder builds IP lists.
type Builder struct {
	client   *http.Client
	logger   Logger
	lookupIP func(host string) ([]net.IP, error)
}

// Logger represents a minimal logger interface.
type Logger interface {
	Debug(s string)
	Info(s string)
	Infof(format string, args ...any)
	Warn(s string)
}

// New returns a new builder of IP lists.
func New(client *http.Client, logger Logger) *Builder {
	return &Builder{
		client:   client,
		logger:   logger,
		lookupIP: net.LookupIP,
	}
}
