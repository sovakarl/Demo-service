package app

import (
	"context"
	"demo-service/config"
	"demo-service/internal/cache"
	"demo-service/internal/cache/memory"
	"demo-service/internal/repository"
	"demo-service/internal/repository/postgres"
	"demo-service/internal/service"
	"demo-service/internal/transport/rest"
	"demo-service/internal/transport/rest/handler/order"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server *http.Server
	db     repository.Repository
	cache  cache.Cache
	logger *slog.Logger
	closer io.Closer
}

func NewApp(cnf *config.Config, logger *slog.Logger) (*App, error) {
	dbConfig := postgres.Config{
		DbName:   cnf.DataBase.DbName,
		Host:     cnf.DataBase.Host,
		Port:     cnf.DataBase.Port,
		Password: cnf.DataBase.Password,
		User:     cnf.DataBase.User,
	}

	// closer := newCloser()

	logger.Info("connect to DataBase...")
	db, err := postgres.NewConnect(dbConfig)
	if err != nil {
		logger.Error("error conect to DataBase", "error", err)
		return nil, err
	}
	// closer.add(db)

	//Конфиг для сервиса
	serviceConfig := service.Config{
		CacheWarmUpLimit: cnf.App.CacheWarmUpLimit,
	}

	cache := memory.NewCache(time.Second*15, time.Second*30, 10, logger)
	// closer.add(cache)

	service := service.NewService(db, cache, serviceConfig, logger)
	Orderhandler := order.NewOrderHandler(service)
	mux := rest.NewOrderRouter(Orderhandler)

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

	// closer.add(serverFuncClose)

	app := &App{
		// closer: closer,
		server: server,
		db:     db,
		cache:  cache,
		logger: logger,
	}
	return app, nil
}

func (a *App) Run() error {
	a.logger.Info("Starting server on", "addr", a.server.Addr)
	return a.server.ListenAndServe()

}

func (a *App) Shutdown(ctx context.Context) error {
	// return a.closer.Close()
	return nil
}
