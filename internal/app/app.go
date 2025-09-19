package app

import (
	"context"
	"demo-service/config"
	"demo-service/internal/cache"
	"demo-service/internal/cache/memory"
	"demo-service/internal/repository"
	"demo-service/internal/repository/postgres"
	"demo-service/internal/service"
	"demo-service/internal/transport/consumer"
	"demo-service/internal/transport/rest"
	"demo-service/internal/transport/rest/handler/order"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server *http.Server
	db     repository.Repository
	cache  cache.Cache
	logger *slog.Logger
	closer *appCloser
	broker *consumer.Kafka
}

func NewApp(cnf *config.Config, logger *slog.Logger) (*App, error) {
	dbConfig := postgres.Config{
		DbName:   cnf.DataBase.DbName,
		Host:     cnf.DataBase.Host,
		Port:     cnf.DataBase.Port,
		Password: cnf.DataBase.Password,
		User:     cnf.DataBase.User,
	}
	closer := newCloser(logger)

	logger.Info("connect to DataBase...")
	db, err := postgres.NewConnect(dbConfig, logger)
	if err != nil {
		return nil, err
	}

	closer.add(db.Close, "databases closed")

	//Конфиг для сервиса
	serviceConfig := service.Config{
		CacheWarmUpLimit: 5,
	}

	cache := memory.NewCache(time.Second*15, time.Second*30, 3, logger)
	closer.add(cache.Close, "cache closed")

	service := service.NewService(db, cache, serviceConfig, logger)
	Orderhandler := order.NewOrderHandler(service)
	mux := rest.NewOrderRouter(Orderhandler)

	brokerConf := consumer.Config{
		Host:    cnf.Consumer.Host,
		Port:    cnf.Consumer.Port,
		GroupID: cnf.Consumer.GroupID,
	}
	broker, err := consumer.NewKafka(brokerConf, service.SaveOrder)
	if err != nil {
		return nil, err
	}
	closer.add(broker.Close, "kafka connecnt closed")

	//Конфиг для запуска сервака
	appConfig := config.App{
		Host: cnf.App.Host,
		Port: cnf.App.Port,
	}

	addr := fmt.Sprintf("%s:%v", appConfig.Host, appConfig.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	app := &App{
		closer: closer,
		server: server,
		db:     db,
		cache:  cache,
		logger: logger,
		broker: broker,
	}

	closer.addWithContext(server.Shutdown, "server stopped")
	return app, nil
}

func (a *App) Run() error {
	// chSignal := make(chan error)
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	a.logger.Info("Starting server on", "addr", a.server.Addr)
	// 	chSignal <- a.server.ListenAndServe()
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	chSignal <- a.broker.Close()
	// }()

	// wg.Wait()
	// return nil
	a.broker.Start()
	a.logger.Info("Starting server on", "addr", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.closer.Close(ctx)
}
