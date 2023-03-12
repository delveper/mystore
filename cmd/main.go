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
	log := lgr.New()

	// Database connection
	db, err := psql.Connect()
	if err != nil {
		log.Fatalf("Failed establishing database connection: %v", err)
	}

	_ = db

	prod := rest.NewProduct(log)

	mux := rest.NewRouter(
		prod.HandleEndpoint,
		// to be continue...
	)

	// Run server
	srv, err := rest.NewServer(mux)
	if err != nil {
		log.Fatalf("Failed creating server: %v", err)
	}

	log.Infof("Server is startign on port %v", os.Getenv("SRV_PORT"))

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed running server: %v", err)
	}

}
