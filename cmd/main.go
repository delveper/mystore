package main

import (
	"log"
	"os"

	"github.com/delveper/mystore/app/interactors"
	repo "github.com/delveper/mystore/app/repository/psql"
	"github.com/delveper/mystore/app/transport/rest"
	"github.com/delveper/mystore/lib/env"
	"github.com/delveper/mystore/lib/lgr"
)

func main() {
	// Load .env
	if err := env.LoadVars(); err != nil {
		log.Fatal(err)
	}

	// Logger
	logger := lgr.New()

	// Database connection
	conn, err := repo.Connect()
	if err != nil {
		logger.Fatalf("Failed establishing database connection: %v", err)
	}

	// Repo
	prodRepo := repo.NewProduct(conn)

	// Interactor
	prodInter := interactors.NewProductInteractor(prodRepo, logger)

	// Handler
	prodREST := rest.NewProduct(prodInter, logger)

	// Mux
	mux := rest.NewMux(prodREST.Route) // routes all endpoints

	// Middleware
	hdl := rest.ChainMiddlewares(mux,
		// top to bottom execution order
		rest.WithLogRequest(logger),
		rest.WithCORS,
		rest.WithJSON,
		rest.WithAuth,
		rest.WithoutPanic(logger),
	)

	// Run server
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
