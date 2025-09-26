package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	_ "github.com/breml/rootcerts"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/qdm12/goservices"
	"github.com/qdm12/goservices/hooks"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/reader/sources/env"
	"github.com/qdm12/gosplash"
	"github.com/qdm12/log"
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
	logger := log.New()
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

//nolint:funlen,cyclop
func _main(ctx context.Context, args []string, logger log.LoggerInterface,
	reader *reader.Reader,
) (err error) {
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		return client.Query(ctx)
	}

	splashSettings := gosplash.MakeLines(gosplash.Settings{
		User:       "qdm12",
		Repository: "updated",
		Emails:     []string{"quentin.mcgaw@gmail.com"},
	})
	for _, line := range splashSettings {
		fmt.Println(line)
	}

	logLevel, err := getLogLevel(reader)
	if err != nil {
		return fmt.Errorf("getting log level: %w", err)
	}
	logger.Patch(log.SetLevel(logLevel))

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

	const healthServerAddr = "127.0.0.1:9999"
	healthServer, setHealthErr, err := health.NewServer(
		healthServerAddr,
		logger.New(log.SetComponent("healthcheck server")),
	)
	if err != nil {
		return fmt.Errorf("creating health server: %w", err)
	}

	runner := run.New(allSettings, logger, shoutrrrSender, shoutrrrParams, setHealthErr)

	sequence, err := goservices.NewSequence(goservices.SequenceSettings{
		ServicesStart: []goservices.Service{healthServer, runner},
		ServicesStop:  []goservices.Service{runner, healthServer},
		Hooks:         hooks.NewWithLog(logger),
	})
	if err != nil {
		return fmt.Errorf("creating sequence of services: %w", err)
	}

	runError, err := sequence.Start(ctx)
	if err != nil {
		return fmt.Errorf("starting services: %w", err)
	}

	select {
	case <-ctx.Done():
		err = sequence.Stop()
		if err != nil {
			return fmt.Errorf("stopping services: %w", err)
		}
		return nil
	case err = <-runError:
		return err
	}
}

func getLogLevel(reader *reader.Reader) (log.Level, error) {
	var settings settings.Log
	settings.Read(reader)
	settings.SetDefaults()
	err := settings.Validate()
	if err != nil {
		return 0, fmt.Errorf("validating log settings: %w", err)
	}
	return log.ParseLevel(settings.Level)
}
