package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	httpapi "github.com/duckvoid/yago-mart/internal/api/http"
	"github.com/duckvoid/yago-mart/internal/config"
	"github.com/duckvoid/yago-mart/internal/logger"
)

type Server struct {
	cfg *config.ServerConfig
	srv *http.Server
}

func New(cfg *config.ServerConfig, handlers httpapi.Handlers) *Server {

	apiRouter := httpapi.NewAPIRouter(handlers)

	return &Server{
		cfg: cfg,
		srv: &http.Server{Addr: cfg.Address, Handler: apiRouter},
	}
}

func (s *Server) Run(ctx context.Context) error {
	logger.Log.Info("Starting server on", slog.String("address", s.srv.Addr))

	listener, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		logger.Log.Error(err.Error())
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
		logger.Log.Error(err.Error())
		return err
	case <-ctx.Done():
		if err := s.srv.Shutdown(ctx); err != nil {
			return err
		}
		logger.Log.Info("Server shutting down")
	}

	return nil
}
