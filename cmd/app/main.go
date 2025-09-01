package main

import (
	"demo-service/config"
	memory "demo-service/internal/cache/in_memory"
	"demo-service/internal/repository/postgres"
	"demo-service/internal/service"
	"demo-service/internal/transport/rest"
	"log"
)

func main() {
	cnf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbConfig := postgres.Config{
		DbName:   cnf.DataBase.DbName,
		Host:     cnf.DataBase.Host,
		Port:     cnf.DataBase.Port,
		Password: cnf.DataBase.Password,
		User:     cnf.DataBase.User,
	}

	db, err := postgres.NewConnect(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cache := memory.NewCache()
	service := service.NewService(db, cache)
	// handler := rest.NewHandler(service)

}
