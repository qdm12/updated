package run

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/updated/internal/constants"
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
	gotify           admin.Gotify
	client           network.Client
	fileManager      files.FileManager
	ipsBuilder       ips.Builder
	hostnamesBuilder hostnames.Builder
	dnscrypto        dnscrypto.DNSCrypto
}

func New(settings settings.Settings, client network.Client,
	logger logging.Logger, gotify admin.Gotify) Runner {
	return &runner{
		settings:         settings,
		logger:           logger,
		gotify:           gotify,
		client:           client,
		ipsBuilder:       ips.NewBuilder(client, logger),
		hostnamesBuilder: hostnames.NewBuilder(client, logger),
		dnscrypto:        dnscrypto.NewDNSCrypto(client, settings.HexSums.NamedRootMD5, settings.HexSums.RootAnchorsSHA256),
		fileManager:      files.NewFileManager(),
	}
}

func (r *runner) Run(ctx context.Context, wg *sync.WaitGroup, period time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	if err := r.singleRun(ctx); err != nil {
		r.gotify.NotifyAndLog(constants.ProgramName, logging.ErrorLevel,
			r.logger, err.Error())
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.singleRun(ctx); err != nil {
				r.gotify.NotifyAndLog(constants.ProgramName, logging.ErrorLevel,
					r.logger, err.Error())
			}
		}
	}
}

func (r *runner) singleRun(ctx context.Context) (err error) {
	tStart := time.Now()
	defer func() {
		executionTime := time.Since(tStart)
		r.logger.Info("overall execution took %s", executionTime)
		r.logger.Info("sleeping for %s", r.settings.Period-executionTime)
	}()
	var gitClient git.Client
	gitSettings := r.settings.Git
	if gitSettings != nil {
		// Setup Git repository
		gitClient, err = git.NewClient(
			gitSettings.SSHKnownHosts,
			gitSettings.SSHKey,
			gitSettings.SSHKeyPassphrase,
			gitSettings.GitURL,
			r.settings.OutputDir)
		if err != nil {
			return err
		}
		if err := gitClient.Pull(); err != nil {
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
		return fmt.Errorf(strings.Join(errorMessages, "; "))
	}
	if gitClient != nil {
		message := fmt.Sprintf("Update of %s", time.Now().Format("2006-01-02"))
		return gitClient.UploadAllChanges(message)
	}
	return nil
}
