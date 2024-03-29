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
	ctx := context.Background()
	args := os.Args
	os.Exit(_main(ctx, args))
}

func _main(ctx context.Context, args []string) (exitCode int) {
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
	envParams := libparams.New()
	level, err := envParams.LogLevel("LOG_LEVEL", libparams.Default("info"))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	logger := logging.New(logging.Settings{Level: level})
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
		logger.Error(err.Error())
		return 1
	}
	client := &http.Client{
		Timeout: HTTPTimeout,
	}
	shoutrrrSender, shoutrrrParams, err := setupShoutrrr(envParams, logger)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	getter := params.NewGetter(envParams)
	allSettings, err := settings.Get(getter)
	if err != nil {
		logger.Error(err.Error())
		return 1
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
		logger.Warn("Caught OS signal " + fmt.Sprint(signal) + ", shutting down")
		cancel()
	}
	wg.Wait()
	return 1
}

func setupShoutrrr(envParams libparams.Interface, logger logging.Logger) (
	sender *router.ServiceRouter, params *types.Params, err error) {
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
