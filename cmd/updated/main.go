package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network/connectivity"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/internal/constants"
	"github.com/qdm12/updated/internal/funcs"
	"github.com/qdm12/updated/internal/health"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/internal/run"
	"github.com/qdm12/updated/internal/settings"
)

func main() {
	ctx := context.Background()
	args := os.Args
	osOpenFile := os.OpenFile
	os.Exit(_main(ctx, args, osOpenFile))
}

func _main(ctx context.Context, args []string, osOpenFile funcs.OSOpenFile) (exitCode int) {
	if health.IsClientMode(args) {
		// Running the program in a separate instance through the Docker
		// built-in healthcheck, in an ephemeral fashion to query the
		// long running instance of the program about its status
		client := health.NewClient()
		if err := client.Query(ctx); err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	envParams := libparams.NewEnv()
	level, err := envParams.LogLevel("LOG_LEVEL", libparams.Default("info"))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	logger := logging.NewParent(logging.Settings{Level: level})
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Print(`
#####################################
############## Updated ##############
########## by Quentin McGaw #########
##### github.com/qdm12/updated ######
#####################################
`)
	HTTPTimeout, err := envParams.Duration("HTTP_TIMEOUT", libparams.Default("10s"))
	if err != nil {
		logger.Error(err)
		return 1
	}
	client := &http.Client{
		Timeout: HTTPTimeout,
	}
	gotify, err := setupGotify(envParams)
	if err != nil {
		logger.Error(err)
		return 1
	}
	getter := params.NewGetter(envParams, osOpenFile)
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
		connChecker := connectivity.NewConnectivity(net.DefaultResolver, client)
		errs := connChecker.Checks(ctx, "github.com")
		for _, err := range errs {
			logger.Warn(err)
		}
	}()

	const healthServerAddr = "127.0.0.1:9999"
	healthServer := health.NewServer(
		healthServerAddr,
		logger.NewChild(logging.Settings{Prefix: "healthcheck server: "}),
	)
	wg.Add(1)
	go healthServer.Run(ctx, wg)

	runner := run.New(allSettings, client, osOpenFile, logger, gotify, healthServer.SetHealthErr)
	// TODO context and in its own goroutine
	logger.Info("Program started")
	if gotify != nil {
		err := gotify.Notify(constants.ProgramName, 1, "Program started")
		if err != nil {
			logger.Error(err.Error())
		}
	}
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

func setupGotify(envParams libparams.Env) (admin.Gotify, error) {
	URL, err := envParams.URL("GOTIFY_URL")
	if err != nil {
		return nil, err
	} else if URL == nil {
		return nil, nil
	}
	token, err := envParams.Get("GOTIFY_TOKEN", libparams.Compulsory())
	if err != nil {
		return nil, err
	}
	return admin.NewGotify(*URL, token, &http.Client{Timeout: time.Second}), nil
}
