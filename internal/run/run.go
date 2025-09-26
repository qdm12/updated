// Package run contains the main update loop runner.
package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/qdm12/updated/internal/settings"
	"github.com/qdm12/updated/pkg/dnscrypto"
	"github.com/qdm12/updated/pkg/hostnames"
	"github.com/qdm12/updated/pkg/ips"
)

// Runner runs the main update loop.
type Runner struct {
	settings         settings.Settings
	logger           Logger
	shoutrrrSender   *router.ServiceRouter
	shoutrrrParams   *types.Params
	ipsBuilder       *ips.Builder
	hostnamesBuilder *hostnames.Builder
	dnscrypto        *dnscrypto.DNSCrypto
	setHealthErr     func(err error)

	// State
	cancel context.CancelFunc
	done   <-chan struct{}
}

// Logger represents a minimal logger interface.
type Logger interface {
	Debug(s string)
	Debugf(format string, args ...any)
	Info(s string)
	Infof(format string, args ...any)
	Warn(s string)
	Error(s string)
}

// New creates a new [Runner] implementing the goservices.Service interface.
func New(settings settings.Settings, logger Logger,
	shoutrrrSender *router.ServiceRouter, shoutrrrParams *types.Params,
	setHealthErr func(err error),
) *Runner {
	client := &http.Client{
		Timeout: settings.HTTPTimeout,
	}
	return &Runner{
		settings:         settings,
		logger:           logger,
		shoutrrrSender:   shoutrrrSender,
		shoutrrrParams:   shoutrrrParams,
		ipsBuilder:       ips.New(client, logger),
		hostnamesBuilder: hostnames.New(client, logger),
		dnscrypto:        dnscrypto.New(client, *settings.HexSums.NamedRootMD5, settings.HexSums.RootAnchorsSHA256),
		setHealthErr:     setHealthErr,
	}
}

func (r *Runner) String() string {
	return "update runner"
}

// Start starts the runner.
func (r *Runner) Start(_ context.Context) (runErr <-chan error, err error) {
	done := make(chan struct{})
	r.done = done
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	ready := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(r.settings.Period)
		defer ticker.Stop()
		close(ready)
		err = r.singleRun(ctx)
		if ctx.Err() != nil {
			return
		}
		r.handleRunError(err)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = r.singleRun(ctx)
				if ctx.Err() != nil {
					return
				}
				r.handleRunError(err)
			}
		}
	}()
	<-ready

	r.shoutrrrSend(r.String() + " started")

	return nil, nil //nolint:nilnil
}

// Stop stops the runner.
func (r *Runner) Stop() error {
	r.cancel()
	<-r.done
	return nil
}

var errEncountered = errors.New("at least one error encountered")

func (r *Runner) singleRun(ctx context.Context) (err error) {
	tStart := time.Now()
	defer func() {
		executionTime := time.Since(tStart)
		r.logger.Infof("overall execution took %s", executionTime)
		r.logger.Infof("sleeping for %s", r.settings.Period-executionTime)
	}()
	gitUploader, err := setupGit(ctx, r.settings, r.logger)
	if err != nil {
		return fmt.Errorf("setting up Git: %w", err)
	}

	chError := make(chan error)
	go func() {
		chError <- r.buildNamedRoot(ctx)
	}()
	go func() {
		chError <- r.buildRootAnchorsAndKeys(ctx)
	}()
	go func() {
		chError <- r.buildMalicious(ctx)
	}()
	go func() {
		chError <- r.buildAds(ctx)
	}()
	go func() {
		chError <- r.buildSurveillance(ctx)
	}()
	var errorMessages []string
	for range 5 {
		err := <-chError
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	close(chError)
	if errorMessages != nil {
		return fmt.Errorf("%w: %s", errEncountered, strings.Join(errorMessages, "; "))
	}
	message := "Update of " + time.Now().Format("2006-01-02")
	err = gitUploader.UploadAllChanges(ctx, message)
	if err != nil {
		return fmt.Errorf("uploading changes: %w", err)
	}
	return nil
}

func (r *Runner) handleRunError(err error) {
	if err == nil {
		r.setHealthErr(nil)
		return
	}
	r.setHealthErr(err)
	r.logger.Error(err.Error())
	r.shoutrrrSend(err.Error())
}

func (r *Runner) shoutrrrSend(message string) {
	errs := r.shoutrrrSender.Send(message, r.shoutrrrParams)
	for _, err := range errs {
		if err != nil {
			r.logger.Warn(err.Error())
		}
	}
}
