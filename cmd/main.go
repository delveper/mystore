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
	prod := rest.NewProduct(nil, logger)

	mux := rest.NewMux(
		prod.Route,
		// to be continue...
	)

	// Chain Middleware
	hdl := rest.ChainMiddlewares(mux, // order matters
		rest.WithLogRequest(logger), // It makes sense to log the request before any other middleware runs
		rest.WithCORS,               // should run early on, so that the appropriate headers can be set to allow cross-origin requests
		rest.WithAuth,               // should run after WithCORS, since authentication headers should be allowed by CORS
		rest.WithJSON,               // run last, so that it can wrap the final response in JSON format
	)

	// Serve server
	srv, err := rest.NewServer(hdl, logger)
	if err != nil {
		logger.Fatalf("Failed creating server: %v", err)
	}

	port := os.Getenv("SRV_PORT")
	logger.Infof("Server is starting on port: %v", port)

	if err := srv.Serve(); err != nil {
		logger.Errorf("Failed running server: %v", err)
	}
}
