package ips

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

var _ Interface = (*Builder)(nil)

type Interface interface {
	BuildMalicious(ctx context.Context) (IPs []string, err error)
	BuildIPsFromHostnames(hostnames []string) (IPs []string)
	CleanIPs(IPs []string) (cleanIPs []string, removedCount int, warnings []string)
}

type Builder struct {
	client   *http.Client
	logger   logging.Logger
	verifier verification.Verifier
	lookupIP func(host string) ([]net.IP, error)
}

func New(client *http.Client, logger logging.Logger) *Builder {
	return &Builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
		lookupIP: net.LookupIP,
	}
}
