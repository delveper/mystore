package rest

import (
	"context"
	"fmt"
	"os/signal"

	"net/http"
	"os"
	"time"

	"github.com/delveper/mystore/lib/lgr"
)

const defaultTimeout = 15 * time.Second

type Server struct {
	server *http.Server
	logger lgr.Logger
}

func NewServer(hdl http.Handler, logger lgr.Logger) (Server, error) {
	addr := os.Getenv("SRV_HOST") + ":" + os.Getenv("SRV_PORT")

	readTimeout, err := time.ParseDuration(os.Getenv("SRV_READ_TIMEOUT"))
	if err != nil {
		return Server{}, fmt.Errorf("parsing read timeout: %w", err)
	}

	writeTimeout, err := time.ParseDuration(os.Getenv("SRV_WRITE_TIMEOUT"))
	if err != nil {
		return Server{}, fmt.Errorf("parsing write timeout: %w", err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("SRV_IDLE_TIMEOUT"))
	if err != nil {
		return Server{}, fmt.Errorf("parsing idle timeout: %w", err)
	}

	srv := http.Server{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      hdl,
	}

	return Server{
		server: &srv,
		logger: logger,
	}, nil
}

func (s *Server) Serve() (err error) {
	go func() {
		if e := s.server.ListenAndServe(); e != nil {
			err = fmt.Errorf("serving: %w", e)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if e := s.server.Shutdown(ctx); e != nil {
		err = fmt.Errorf("shutting down server: %w", e)
	}

	s.logger.Info("Shutting down gracefully.")

	return
}
