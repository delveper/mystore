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

// defaultTimeout is the default duration for shutting down the server gracefully.
const defaultTimeout = 15 * time.Second

// Server defines a RESTful server.
type Server struct {
	server *http.Server
	logger lgr.Logger
}

// NewServer creates a new instance of Server with specified http.Handler and logger.
// It reads environment variables to set up the server configuration.
func NewServer(hdl http.Handler, logger lgr.Logger) (Server, error) {
	// Get server configuration from environment variables.
	addr := os.Getenv("SRV_HOST") + ":" + os.Getenv("SRV_PORT")

	readTimeout, err := time.ParseDuration(os.Getenv("SRV_READ_TIMEOUT"))
	if err != nil {
		return Server{}, fmt.Errorf("parsing Read timeout: %w", err)
	}

	writeTimeout, err := time.ParseDuration(os.Getenv("SRV_WRITE_TIMEOUT"))
	if err != nil {
		return Server{}, fmt.Errorf("parsing write timeout: %w", err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("SRV_IDLE_TIMEOUT"))
	if err != nil {
		return Server{}, fmt.Errorf("parsing idle timeout: %w", err)
	}

	// Create http.Server instance.
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

// Serve runs the server until it receives an interrupt signal (Ctrl+C).
// It gracefully shuts down the server and logs the event.
func (s *Server) Serve() (err error) {
	// Start the server in a new goroutine.
	go func() {
		if e := s.server.ListenAndServe(); e != nil {
			err = fmt.Errorf("serving: %w", e)
		}
	}()

	// Wait for an interrupt signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a new context with a default timeout.
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if e := s.server.Shutdown(ctx); e != nil {
		err = fmt.Errorf("shutting down server: %w", e)
	}

	s.logger.Info("Shutting down gracefully.")

	return
}
