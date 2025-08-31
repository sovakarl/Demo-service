package main

import (
	"demo-service/config"
	"demo-service/storage/postgres"
	"log"
)

func main() {
	cnf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	db, err := postgres.NewConnect(postgres.Config(cnf.DataBase))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
