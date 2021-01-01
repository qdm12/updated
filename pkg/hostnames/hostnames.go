package hostnames

import (
	"context"
	"net/http"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

type Builder interface {
	BuildSurveillance(ctx context.Context) (hostnames []string, err error)
	BuildMalicious(ctx context.Context) (hostnames []string, err error)
	BuildAds(ctx context.Context) (hostnames []string, err error)
}

type builder struct {
	client   *http.Client
	logger   logging.Logger
	verifier verification.Verifier
}

func NewBuilder(client *http.Client, logger logging.Logger) Builder {
	return &builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
	}
}
