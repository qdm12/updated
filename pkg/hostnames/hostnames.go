package hostnames

import (
	"context"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
)

type Builder interface {
	BuildSurveillance(ctx context.Context) (hostnames []string, err error)
	BuildMalicious(ctx context.Context) (hostnames []string, err error)
	BuildAds(ctx context.Context) (hostnames []string, err error)
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
