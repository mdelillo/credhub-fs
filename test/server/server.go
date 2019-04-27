package server

import (
	"context"
	"net/http"
	"time"
)

type Server interface {
	Start() error
	Shutdown() error
}

type server struct {
	listenAddr string
	certPath   string
	keyPath    string
	handler    http.Handler
	httpServer *http.Server
}

func NewServer(listenAddr, certPath, keyPath string, handler http.Handler) Server {
	return &server{
		listenAddr: listenAddr,
		certPath:   certPath,
		keyPath:    keyPath,
		handler:    handler,
	}
}

func (s *server) Start() error {
	s.httpServer = &http.Server{
		Addr:    s.listenAddr,
		Handler: s.handler,
	}

	if err := s.httpServer.ListenAndServeTLS(s.certPath, s.keyPath); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) Shutdown() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}

		s.httpServer = nil
	}
	return nil
}
