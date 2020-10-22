package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/network/connectivity"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/internal/constants"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/internal/run"
	"github.com/qdm12/updated/internal/settings"
)

func main() {
	ctx := context.Background()
	os.Exit(_main(ctx))
}

func _main(ctx context.Context) (exitCode int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	envParams := libparams.NewEnvParams()
	encoding, level, err := envParams.GetLoggerConfig()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	logger, err := logging.NewLogger(encoding, level)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if healthcheck.Mode(os.Args) {
		if err := healthcheck.Query(ctx); err != nil {
			logger.Error(err)
			return 1
		}
		return 0
	}
	fmt.Print(`
#####################################
############## Updated ##############
########## by Quentin McGaw #########
##### github.com/qdm12/updated ######
#####################################
`)
	HTTPTimeout, err := envParams.GetHTTPTimeout(libparams.Default("3s"))
	if err != nil {
		logger.Error(err)
		return 1
	}
	client := network.NewClient(HTTPTimeout)
	gotify, err := setupGotify(envParams)
	if err != nil {
		logger.Error(err)
		return 1
	}
	getter := params.NewGetter(envParams)
	allSettings, err := settings.Get(getter)
	if err != nil {
		logger.Error(err)
		return 1
	}
	logger.Info(allSettings.String())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		errs := connectivity.NewConnectivity(HTTPTimeout).Checks(ctx, "github.com")
		for _, err := range errs {
			logger.Warn(err)
		}
	}()
	runner := run.New(allSettings, client, logger, gotify)
	// TODO context and in its own goroutine
	gotify.NotifyAndLog(constants.ProgramName, logging.InfoLevel, logger, "Program started")
	wg.Add(1)
	go runner.Run(ctx, wg, allSettings.Period)
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	select {
	case <-ctx.Done():
		logger.Warn("context canceled, shutting down")
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
		cancel()
	}
	wg.Wait()
	return 1
}

func setupGotify(envParams libparams.EnvParams) (admin.Gotify, error) {
	URL, err := envParams.GetGotifyURL()
	if err != nil {
		return nil, err
	} else if URL == nil {
		return nil, nil
	}
	token, err := envParams.GetGotifyToken()
	if err != nil {
		return nil, err
	}
	return admin.NewGotify(*URL, token, &http.Client{Timeout: time.Second}), nil
}
