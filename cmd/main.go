package main

import (
	"log"

	"github.com/delveper/mystore/app/repository"
	"github.com/delveper/mystore/lib/env"
	"github.com/delveper/mystore/lib/lgr"
)

func main() {
	// Load .env to environment
	if err := env.LoadVars(); err != nil {
		log.Fatal(err)
	}

	// Logger
	log := lgr.New()

	// Database connection
	db, err := repository.Connect()
	if err != nil {
		log.Fatalf("Failed establishing database connection: %v", err)
	}

	_ = db

}
