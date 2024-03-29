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

type Runner interface {
	Run(ctx context.Context, wg *sync.WaitGroup, period time.Duration)
}

type runner struct {
	settings         settings.Settings
	logger           logging.Logger
	shoutrrrSender   *router.ServiceRouter
	shoutrrrParams   *types.Params
	ipsBuilder       ips.Interface
	hostnamesBuilder hostnames.Interface
	dnscrypto        dnscrypto.Interface
	setHealthErr     func(err error)
}

func New(settings settings.Settings, client *http.Client,
	logger logging.Logger, shoutrrrSender *router.ServiceRouter, shoutrrrParams *types.Params,
	setHealthErr func(err error)) Runner {
	return &runner{
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

func (r *runner) Run(ctx context.Context, wg *sync.WaitGroup, period time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	if err := r.singleRun(ctx); err != nil {
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
			if err := r.singleRun(ctx); err != nil {
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

func (r *runner) singleRun(ctx context.Context) (err error) {
	tStart := time.Now()
	defer func() {
		executionTime := time.Since(tStart)
		r.logger.Info(fmt.Sprintf("overall execution took %s", executionTime))
		r.logger.Info(fmt.Sprintf("sleeping for %s", r.settings.Period-executionTime))
	}()
	var gitClient git.Interface
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
		if err := gitClient.Pull(ctx); err != nil {
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
	for N := 0; N < 5; N++ {
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
		message := fmt.Sprintf("Update of %s", time.Now().Format("2006-01-02"))
		if err := gitClient.UploadAllChanges(ctx, message); err != nil {
			return err
		}
		r.logger.Info("Committed to Git: " + message)
	}
	return nil
}
