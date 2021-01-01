package ips

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

type Builder interface {
	BuildMalicious(ctx context.Context) (IPs []string, err error)
	BuildIPsFromHostnames(hostnames []string) (IPs []string)
	CleanIPs(IPs []string) (cleanIPs []string, removedCount int, warnings []string)
}

type builder struct {
	client   *http.Client
	logger   logging.Logger
	verifier verification.Verifier
	lookupIP func(host string) ([]net.IP, error)
}

func NewBuilder(client *http.Client, logger logging.Logger) Builder {
	return &builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
		lookupIP: net.LookupIP,
	}
}
