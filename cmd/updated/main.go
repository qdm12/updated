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
	_ "time/tzdata"

	_ "github.com/breml/rootcerts"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/reader/sources/env"
	"github.com/qdm12/updated/internal/health"
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
	reader := reader.New(reader.Settings{
		Sources: []reader.Source{env.New(env.Settings{})},
	})

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, args, logger, reader)
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

//nolint:funlen
func _main(ctx context.Context, args []string, logger logging.ParentLogger,
	reader *reader.Reader,
) (err error) {
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		return client.Query(ctx)
	}

	fmt.Print(`
#####################################
############## Updated ##############
########## by Quentin McGaw #########
##### github.com/qdm12/updated ######
#####################################
`)

	logLevel, err := getLogLevel(reader)
	if err != nil {
		return fmt.Errorf("getting log level: %w", err)
	}
	logger.PatchLevel(logLevel)

	var allSettings settings.Settings
	err = allSettings.Read(reader)
	if err != nil {
		return fmt.Errorf("reading settings: %w", err)
	}
	allSettings.SetDefaults()
	err = allSettings.Validate()
	if err != nil {
		return fmt.Errorf("validating settings: %w", err)
	}
	logger.Info(allSettings.String())

	shoutrrrSender, err := shoutrrr.CreateSender(allSettings.Shoutrrr.ServiceURLs...)
	if err != nil {
		return fmt.Errorf("setting up Shoutrrr: %w", err)
	}
	shoutrrrParams := &types.Params{}
	shoutrrrParams.SetTitle("Updated")

	wg := &sync.WaitGroup{}

	const healthServerAddr = "127.0.0.1:9999"
	healthServer := health.NewServer(
		healthServerAddr,
		logger.NewChild(logging.Settings{Prefix: "healthcheck server: "}),
	)
	wg.Add(1)
	go healthServer.Run(ctx, wg)

	client := &http.Client{
		Timeout: allSettings.HTTPTimeout,
	}
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

func getLogLevel(reader *reader.Reader) (logging.Level, error) {
	var settings settings.Log
	settings.Read(reader)
	settings.SetDefaults()
	err := settings.Validate()
	if err != nil {
		return 0, fmt.Errorf("validating log settings: %w", err)
	}
	levels := []logging.Level{logging.LevelDebug, logging.LevelInfo, logging.LevelWarn, logging.LevelError}
	for _, level := range levels {
		if settings.Level == level.String() {
			return level, nil
		}
	}
	panic("log level not recognized: " + settings.Level) // should be validated earlier
}
