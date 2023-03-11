package main

import (
	"log"

	"github.com/delveper/mystore/app/repo"
	"github.com/delveper/mystore/lib/env"
)

func main() {
	if err := env.LoadVars(); err != nil {
		log.Fatal(err)
	}

	db, err := repo.Connect()
	if err != nil {
		log.Fatal(err)
	}

	_ = db
}
