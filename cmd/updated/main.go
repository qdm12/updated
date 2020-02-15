package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/qdm12/updated/internal/env"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/internal/run"
	"github.com/qdm12/updated/internal/settings"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/signals"
)

func main() {
	logger, err := logging.NewLogger(logging.JSONEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	envParams := libparams.NewEnvParams()
	encoding, level, nodeID, err := envParams.GetLoggerConfig()
	if err != nil {
		logger.Error(err)
	} else {
		logger, err = logging.NewLogger(encoding, level, nodeID)
	}
	if healthcheck.Mode(os.Args) {
		if err := healthcheck.Query(); err != nil {
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
	paramsGetter := params.NewParamsGetter(envParams)
	allSettings, err := settings.Get(paramsGetter)
	logger.Info(allSettings.String())
	go signals.WaitForExit(e.ShutdownFromSignal)
	errs := network.NewConnectivity(HTTPTimeout).Checks("github.com")
	for _, err := range errs {
		e.Warn(err)
	}
	runner := run.NewRunner(allSettings, client, logger)
	e.Notify(1, "Program started")
	for {
		go func() {
			err := runner.Run()
			e.CheckError(err)
		}()
		time.Sleep(allSettings.Period)
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
