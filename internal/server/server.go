package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	httpapi "github.com/duckvoid/yago-mart/internal/api/http"
	"github.com/duckvoid/yago-mart/internal/config"
)

type Server struct {
	cfg    *config.ServerConfig
	srv    *http.Server
	logger *slog.Logger
}

func New(cfg *config.ServerConfig, handlers httpapi.Handlers, logger *slog.Logger) *Server {

	apiRouter := httpapi.NewAPIRouter(handlers)

	return &Server{
		cfg:    cfg,
		srv:    &http.Server{Addr: cfg.Address, Handler: apiRouter},
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting server on", slog.String("address", s.srv.Addr))

	listener, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		s.logger.Error("Failed to create listener", "error", err)
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		s.logger.Error("Failed to start server", "error", err)
		return err
	case <-ctx.Done():
		if err := s.srv.Shutdown(ctx); err != nil {
			s.logger.Error("Failed to shutdown server", "error", err)
			return err
		}
		s.logger.Info("Server shut down")
	}

	return nil
}
