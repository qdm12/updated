package env

import (
	"os"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

// Env contains objects and methods necessary to the main function.
// These are created at start and are needed to the top-level
// working management of the program.
type Env interface {
	SetGotify(gotify admin.Gotify)
	SetClient(client network.Client)
	Notify(priority int, args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	CheckError(err error)
	FatalOnError(err error)
	ShutdownFromSignal(signal string) (exitCode int)
	Fatal(args ...interface{})
	Shutdown() (exitCode int)
}

type env struct {
	client network.Client
	logger logging.Logger
	gotify admin.Gotify
}

// NewEnv creates a new Env object
func NewEnv(logger logging.Logger) Env {
	return &env{logger: logger}
}

func (e *env) SetGotify(gotify admin.Gotify) {
	e.gotify = gotify
}

func (e *env) SetClient(client network.Client) {
	e.client = client
}

// Notify sends a notification to the Gotify server.
func (e *env) Notify(priority int, args ...interface{}) {
	if e.gotify != nil {
		err := e.gotify.Notify("Updated", priority, args...)
		if err != nil {
			e.logger.Error(err)
		}
	}
}

// Info logs a message and sends a notification to the Gotify server.
func (e *env) Info(args ...interface{}) {
	e.logger.Info(args...)
	e.Notify(1, args...)
}

// Warn logs a message and sends a notification to the Gotify server.
func (e *env) Warn(args ...interface{}) {
	e.logger.Warn(args...)
	e.Notify(2, args...)
}

// CheckError logs an error and sends a notification to the Gotify server
// if the error is not nil.
func (e *env) CheckError(err error) {
	if err != nil {
		e.logger.Error(err)
		e.Notify(3, err)
	}
}

// FatalOnError calls Fatal if the error is not nil.
func (e *env) FatalOnError(err error) {
	if err != nil {
		e.Fatal(err)
	}
}

// Shutdown cleanly exits the program by closing all connections,
// databases and syncing the loggers.
func (e *env) Shutdown() (exitCode int) {
	defer func() {
		if err := e.logger.Sync(); err != nil {
			exitCode = 99
		}
	}()
	if e.client != nil {
		e.client.Close()
	}
	return 0
}

// ShutdownFromSignal logs a warning, sends a notification to Gotify and shutdowns
// the program cleanly when a OS level signal is received. It should be passed as a
// callback to a function which would catch such signal.
func (e *env) ShutdownFromSignal(signal string) (exitCode int) {
	e.logger.Warn("Program stopped with signal %s", signal)
	e.Notify(1, "Caught OS signal "+signal)
	return e.Shutdown()
}

// Fatal logs an error, sends a notification to Gotify and shutdowns the program.
// It exits the program with an exit code of 1.
func (e *env) Fatal(args ...interface{}) {
	e.logger.Error(args...)
	e.Notify(4, args...)
	e.Shutdown()
	os.Exit(1)
}
