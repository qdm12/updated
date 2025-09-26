// Package health provides a health HTTP server and client.
package health

import (
	"github.com/qdm12/goservices/httpserver"
)

// Logger represents a minimal logger interface.
type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}

// NewServer creates a new health server.
func NewServer(address string, logger Logger) (
	server *httpserver.Server, setErr func(error), err error,
) {
	handler := newHandler(logger)
	setErr = handler.setHealthErr
	server, err = httpserver.New(httpserver.Settings{
		Handler: newHandler(logger),
		Name:    ptrTo("healthcheck"),
		Address: ptrTo(address),
		Logger:  logger,
	})
	if err != nil {
		return nil, nil, err
	}
	return server, setErr, nil
}

func ptrTo[T any](v T) *T { return &v }
