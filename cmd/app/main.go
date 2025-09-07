package main

import (
	"demo-service/config"
	"demo-service/internal/cache/memory"
	"demo-service/internal/repository/postgres"
	"demo-service/internal/service"
	"demo-service/internal/transport/rest"
	"demo-service/internal/transport/rest/handler/order"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Closer interface {
	Close()
}

func main() {
	cnf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// канал для системных вызовов
	sigCh := make(chan os.Signal,1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

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

	cache := memory.NewCache(time.Minute, time.Minute*2)
	service := service.NewService(db, cache)
	Orderhandler := order.NewOrderHandler(service)
	mux := rest.NewOrderRouter(Orderhandler)

	go func() {
		http.ListenAndServe(":8080", mux)
	}()

	<-sigCh
	cleanup(db, cache)
}

func cleanup(objects ...Closer) {
	for _, object := range objects {
		object.Close()
	}
}

// http->router->handler->service->db||cache
