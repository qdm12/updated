package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/network/connectivity"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/internal/env"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/internal/run"
	"github.com/qdm12/updated/internal/settings"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		panic(err)
	}
	envParams := libparams.NewEnvParams()
	encoding, level, err := envParams.GetLoggerConfig()
	if err != nil {
		logger.Error(err)
	} else {
		logger, err = logging.NewLogger(encoding, level)
		if err != nil {
			panic(err)
		}
	}
	if healthcheck.Mode(os.Args) {
		if err := healthcheck.Query(ctx); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	fmt.Print(`
#####################################
############## Updated ##############
########## by Quentin McGaw #########
##### github.com/qdm12/updated ######
#####################################
`)
	e := env.NewEnv(logger)
	HTTPTimeout, err := envParams.GetHTTPTimeout(libparams.Default("3s"))
	e.FatalOnError(err)
	client := network.NewClient(HTTPTimeout)
	e.SetClient(client)
	gotify, err := setupGotify(envParams)
	if err != nil {
		logger.Error(err)
	} else {
		e.SetGotify(gotify)
	}
	getter := params.NewGetter(envParams)
	allSettings, err := settings.Get(getter)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(allSettings.String())
	errs := connectivity.NewConnectivity(HTTPTimeout).Checks(ctx, "github.com")
	for _, err := range errs {
		e.Warn(err)
	}
	runner := run.NewRunner(allSettings, client, logger)
	e.Notify(1, "Program started")
	go func() {
		ticker := time.NewTicker(allSettings.Period)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				err := runner.Run(ctx)
				e.CheckError(err)
			}
		}
	}()
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
