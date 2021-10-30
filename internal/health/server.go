package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	SetHealthErr(err error)
}

type server struct {
	address string
	logger  logging.Logger
	handler *handler
}

func NewServer(address string, logger logging.Logger) Server {
	return &server{
		address: address,
		logger:  logger,
		handler: newHandler(logger),
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		s.logger.Warn("shutting down (context canceled)")
		defer s.logger.Warn("shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
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

func (s *server) SetHealthErr(err error) {
	s.handler.setHealthErr(err)
}
