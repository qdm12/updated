// Package health provides a health HTTP server and client.
package health

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Server represents a health HTTP server.
type Server struct {
	address string
	logger  Logger
	handler *handler
}

// Logger represents a minimal logger interface.
type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}

// NewServer creates a new health server.
func NewServer(address string, logger Logger) *Server {
	return &Server{
		address: address,
		logger:  logger,
		handler: newHandler(logger),
	}
}

// Run starts the HTTP server until the context is done.
func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler, ReadHeaderTimeout: time.Second}
	go func() {
		<-ctx.Done()
		s.logger.Warn("shutting down (context canceled)")
		defer s.logger.Warn("shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Error("failed shutting down: " + err.Error())
		}
	}()
	for ctx.Err() == nil {
		s.logger.Info("listening on " + s.address)
		err := server.ListenAndServe()
		if err != nil && ctx.Err() == nil { // server crashed
			s.logger.Error(err.Error())
			s.logger.Info("restarting")
		}
	}
}

// SetHealthErr sets the health error to be returned by the /health endpoint.
func (s *Server) SetHealthErr(err error) {
	s.handler.setHealthErr(err)
}
