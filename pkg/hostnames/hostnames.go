package hostnames

import (
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
)

type Builder interface {
	BuildSurveillance() (hostnames []string, err error)
	BuildMalicious() (hostnames []string, err error)
	BuildAds() (hostnames []string, err error)
}

type builder struct {
	client   network.Client
	logger   logging.Logger
	verifier verification.Verifier
}

func NewBuilder(client network.Client, logger logging.Logger) Builder {
	return &builder{
		client,
		logger,
		verification.NewVerifier(),
	}
}
