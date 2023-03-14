package main

import (
	"log"
	"os"

	"github.com/delveper/mystore/app/repository/psql"
	"github.com/delveper/mystore/app/transport/rest"
	"github.com/delveper/mystore/lib/env"
	"github.com/delveper/mystore/lib/lgr"
)

func main() {
	// Load .env to environment, in case of running locally
	if err := env.LoadVars(); err != nil {
		log.Fatal(err)
	}

	// Logger
	logger := lgr.New()

	// Database connection
	db, err := psql.Connect()
	if err != nil {
		logger.Fatalf("Failed establishing database connection: %v", err)
	}

	_ = db

	// Mux handlers
	prod := rest.NewProduct(logger)

	mux := rest.NewMux(
		prod.Route,
		// to be continue...
	)

	// Chain Middleware
	hdl := rest.ChainMiddlewares(mux,
		rest.WithLogRequest(logger),
		rest.WithJSON,
	)

	// Serve server
	srv, err := rest.NewServer(hdl, logger)
	if err != nil {
		logger.Fatalf("Failed creating server: %v", err)
	}

	port := os.Getenv("SRV_PORT")
	logger.Infof("Server is starting on port %v", port)

	if err := srv.Serve(); err != nil {
		logger.Errorf("Failed running server: %v", err)
	}
}
