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
	srv *http.Server
	log *lgr.Logger
}

func NewServer(hdl http.Handler, log *lgr.Logger) (*Server, error) {
	addr := os.Getenv("SRV_HOST") + ":" + os.Getenv("SRV_PORT")

	readTimeout, err := time.ParseDuration(os.Getenv("SRV_READ_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("parse read timeout: %w", err)
	}

	writeTimeout, err := time.ParseDuration(os.Getenv("SRV_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("parse write timeout: %w", err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("SRV_IDLE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("parse idle timeout: %w", err)
	}

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      hdl,
	}

	return &Server{
		srv: srv,
		log: log,
	}, nil
}

func (s *Server) Serve() (err error) {
	go func() {
		if e := s.srv.ListenAndServe(); e != nil {
			err = fmt.Errorf("serving: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if e := s.srv.Shutdown(ctx); e != nil {
		err = fmt.Errorf("shutting down server: %v", err)
	}

	s.log.Info("Shutting down gracefully.")

	return
}
