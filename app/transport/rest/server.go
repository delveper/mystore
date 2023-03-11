package rest

import (
	"fmt"

	"net/http"
	"os"
	"time"
)

type Server struct{ *http.Server }

func NewServer(hdl http.Handler) (*Server, error) {
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

	return &Server{srv}, nil
}

func (srv *Server) Run() error {
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("run the server: %w", err)
	}

	return nil
}
