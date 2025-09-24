package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	_ "github.com/breml/rootcerts"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/internal/health"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/internal/run"
	"github.com/qdm12/updated/internal/settings"
)

func main() {
	background := context.Background()
	ctx, cancel := context.WithCancel(background)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	args := os.Args
	logger := logging.New(logging.Settings{})

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, args, logger)
	}()

	// Wait for OS signal or run error
	var err error
	select {
	case receivedSignal := <-signalCh:
		signal.Stop(signalCh)
		fmt.Println("")
		logger.Warn("Caught OS signal " + receivedSignal.String() + ", shutting down")
		cancel()
	case err = <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Error(err.Error())
		cancel()
	}

	// Shutdown timed sequence, and force exit on second OS signal
	const shutdownGracePeriod = 5 * time.Second
	timer := time.NewTimer(shutdownGracePeriod)
	select {
	case shutdownErr := <-errorCh:
		timer.Stop()
		if shutdownErr != nil {
			logger.Warn("Shutdown failed: " + shutdownErr.Error())
			os.Exit(1)
		}

		logger.Info("Shutdown successful")
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	case <-timer.C:
		logger.Warn("Shutdown timed out")
		os.Exit(1)
	}
}

func _main(ctx context.Context, args []string, logger logging.ParentLogger) (err error) {
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		return client.Query(ctx)
	}

	envParams := libparams.New()
	level, err := envParams.LogLevel("LOG_LEVEL", libparams.Default("info"))
	if err != nil {
		return fmt.Errorf("getting log level: %w", err)
	}

	logger.PatchLevel(level)

	fmt.Print(`
#####################################
############## Updated ##############
########## by Quentin McGaw #########
##### github.com/qdm12/updated ######
#####################################
`)
	HTTPTimeout, err := envParams.Duration("HTTP_TIMEOUT", libparams.Default("10s"))
	if err != nil {
		return fmt.Errorf("getting HTTP timeout: %w", err)
	}
	client := &http.Client{
		Timeout: HTTPTimeout,
	}
	shoutrrrSender, shoutrrrParams, err := setupShoutrrr(envParams, logger)
	if err != nil {
		return fmt.Errorf("setting up Shoutrrr: %w", err)
	}

	getter := params.NewGetter(envParams)
	allSettings, err := settings.Get(getter)
	if err != nil {
		return fmt.Errorf("getting settings: %w", err)
	}
	logger.Info(allSettings.String())

	wg := &sync.WaitGroup{}

	const healthServerAddr = "127.0.0.1:9999"
	healthServer := health.NewServer(
		healthServerAddr,
		logger.NewChild(logging.Settings{Prefix: "healthcheck server: "}),
	)
	wg.Add(1)
	go healthServer.Run(ctx, wg)

	runner := run.New(allSettings, client, logger, shoutrrrSender, shoutrrrParams, healthServer.SetHealthErr)
	// TODO context and in its own goroutine
	logger.Info("Program started")
	errs := shoutrrrSender.Send("Program started", shoutrrrParams)
	for _, err := range errs {
		if err != nil {
			logger.Error(err.Error())
		}
	}
	wg.Add(1)
	runner.Run(ctx, wg, allSettings.Period) // this can only exit when context is canceled.
	wg.Wait()
	return nil
}

func setupShoutrrr(envParams libparams.Interface, logger logging.Logger) (
	sender *router.ServiceRouter, params *types.Params, err error,
) {
	shoutrrrURLs, err := envParams.Get("SHOUTRRR_SERVICES", libparams.CaseSensitiveValue())
	if err != nil {
		return nil, nil, err
	}
	var rawURLs []string
	if shoutrrrURLs != "" {
		rawURLs = strings.Split(shoutrrrURLs, ",")
		logger.Info("Using " + strconv.Itoa(len(rawURLs)) + "Shoutrrr service URLs")
	}

	sender, err = shoutrrr.CreateSender(rawURLs...)
	if err != nil {
		return nil, nil, err
	}

	params = &types.Params{}
	params.SetTitle("Updated")

	return sender, params, nil
}
