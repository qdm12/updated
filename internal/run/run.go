// Package run contains the main update loop runner.
package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/updated/internal/settings"
	"github.com/qdm12/updated/pkg/dnscrypto"
	"github.com/qdm12/updated/pkg/git"
	"github.com/qdm12/updated/pkg/hostnames"
	"github.com/qdm12/updated/pkg/ips"
)

// Runner runs the main update loop.
type Runner struct {
	settings         settings.Settings
	logger           logging.Logger
	shoutrrrSender   *router.ServiceRouter
	shoutrrrParams   *types.Params
	ipsBuilder       *ips.Builder
	hostnamesBuilder *hostnames.Builder
	dnscrypto        *dnscrypto.DNSCrypto
	setHealthErr     func(err error)
}

// New creates a new Runner.
func New(settings settings.Settings, client *http.Client,
	logger logging.Logger, shoutrrrSender *router.ServiceRouter, shoutrrrParams *types.Params,
	setHealthErr func(err error),
) *Runner {
	return &Runner{
		settings:         settings,
		logger:           logger,
		shoutrrrSender:   shoutrrrSender,
		shoutrrrParams:   shoutrrrParams,
		ipsBuilder:       ips.New(client, logger),
		hostnamesBuilder: hostnames.New(client, logger),
		dnscrypto:        dnscrypto.New(client, settings.HexSums.NamedRootMD5, settings.HexSums.RootAnchorsSHA256),
		setHealthErr:     setHealthErr,
	}
}

// Run starts the main loop that runs every period duration until the context is done.
func (r *Runner) Run(ctx context.Context, wg *sync.WaitGroup, period time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	err := r.singleRun(ctx)
	if err != nil {
		r.setHealthErr(err)
		r.logger.Error(err.Error())
		errs := r.shoutrrrSender.Send(err.Error(), r.shoutrrrParams)
		for _, err := range errs {
			if err != nil {
				r.logger.Error(err.Error())
			}
		}
	} else {
		r.setHealthErr(nil)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := r.singleRun(ctx)
			if err != nil {
				r.setHealthErr(err)
				r.logger.Error(err.Error())
				errs := r.shoutrrrSender.Send(err.Error(), r.shoutrrrParams)
				for _, err := range errs {
					if err != nil {
						r.logger.Error(err.Error())
					}
				}
			} else {
				r.setHealthErr(nil)
			}
		}
	}
}

var errEncountered = errors.New("at least one error encountered")

func (r *Runner) singleRun(ctx context.Context) (err error) {
	tStart := time.Now()
	defer func() {
		executionTime := time.Since(tStart)
		r.logger.Info(fmt.Sprintf("overall execution took %s", executionTime))
		r.logger.Info(fmt.Sprintf("sleeping for %s", r.settings.Period-executionTime))
	}()
	var gitClient *git.Client
	gitSettings := r.settings.Git
	if gitSettings != nil {
		// Setup Git repository
		gitClient, err = git.New(
			gitSettings.SSHKnownHosts,
			gitSettings.SSHKey,
			gitSettings.SSHKeyPassphrase,
			gitSettings.GitURL,
			r.settings.OutputDir)
		if err != nil {
			return err
		}
		err = gitClient.Pull(ctx)
		if err != nil {
			return err
		}
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
	if gitClient != nil {
		message := "Update of " + time.Now().Format("2006-01-02")
		err = gitClient.UploadAllChanges(ctx, message)
		if err != nil {
			return err
		}
		r.logger.Info("Committed to Git: " + message)
	}
	return nil
}
