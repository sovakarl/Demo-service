package main

import (
	"demo-service/config"
	"demo-service/internal/cache/memory"
	"demo-service/internal/repository/postgres"
	"demo-service/internal/service"
	"demo-service/internal/transport/rest"
	"demo-service/internal/transport/rest/handler/order"
	"demo-service/logger"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Closer interface {
	Close()
}

func main() {

	cnf, err := config.Load()
	if err != nil {
		log.Fatal("error pars config")
	}

	logConf := logger.Config{
		LogLvl: cnf.Log.LogLevel,
	}

	logger := logger.NewLoger(logConf)

	// канал для системных вызовов
	// sigCh := make(chan os.Signal,1)
	// signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// конфиг для БД
	dbConfig := postgres.Config{
		DbName:   cnf.DataBase.DbName,
		Host:     cnf.DataBase.Host,
		Port:     cnf.DataBase.Port,
		Password: cnf.DataBase.Password,
		User:     cnf.DataBase.User,
	}

	logger.Info("connect to DataBase...")
	db, err := postgres.NewConnect(dbConfig)
	if err != nil {
		logger.Error("error conect to DataBase", "error", err)
		return
	}

	//Конфиг для сервиса
	serviceConfig := service.Config{
		CacheWarmUpLimit: cnf.App.CacheWarmUpLimit,
	}

	cache := memory.NewCache(time.Minute, time.Minute*2, 10)
	service := service.NewService(db, cache, serviceConfig, logger)
	Orderhandler := order.NewOrderHandler(service)
	mux := rest.NewOrderRouter(Orderhandler)

	//Конфиг для запуска сервака
	appConfig := config.App{
		Host: cnf.App.Host,
		Port: cnf.App.Port,
	}

	// костыль ебаный с этим адресом
	addr := fmt.Sprintf("%s:%v", appConfig.Host, appConfig.Port)
	http.ListenAndServe(addr, mux)

	//чистое завершение проги
	// <-sigCh
	// cleanup(db, cache)
}

func cleanup(objects ...Closer) {
	for _, object := range objects {
		object.Close()
	}
}

// http->router->handler->service->db||cache
