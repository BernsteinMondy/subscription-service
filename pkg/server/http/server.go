package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"net/http"
)

const (
	_defaultPort         = 8080
	_defaultWriteTimeout = 10 * time.Second
	_defaultReadTimeout  = 10 * time.Second
)

type httpServer struct {
	srv *http.Server
}

func NewServer() *httpServer {
	server := &httpServer{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", _defaultPort),
			WriteTimeout: _defaultWriteTimeout,
			ReadTimeout:  _defaultReadTimeout,
		},
	}

	return server
}

func (s *httpServer) Launch() error {
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: listen and serve: %w", err)
	}

	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("could not shutdown http server: %w", err)
	}

	return nil
}
