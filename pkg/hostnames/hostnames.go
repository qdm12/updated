package hostnames

import (
	"context"
	"net/http"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

var _ Interface = (*Builder)(nil)

type Interface interface {
	BuildSurveillance(ctx context.Context) (hostnames []string, err error)
	BuildMalicious(ctx context.Context) (hostnames []string, err error)
	BuildAds(ctx context.Context) (hostnames []string, err error)
}

type Builder struct {
	client   *http.Client
	logger   logging.Logger
	verifier verification.Verifier
}

func New(client *http.Client, logger logging.Logger) *Builder {
	return &Builder{
		client:   client,
		logger:   logger,
		verifier: verification.NewVerifier(),
	}
}
