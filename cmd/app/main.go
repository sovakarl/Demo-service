package main

import (
	"demo-service/configs"
	"demo-service/storage/postgre"
	"log"
)

func main() {
	appConfig := configs.Load()
	
	db,err:=postgre.Connect(cnf.DBConfig)
	if err!=nil{
		log.Fatal(err)
	}
}
