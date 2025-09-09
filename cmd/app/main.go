package main

import (
	"context"
	"demo-service/config"
	"demo-service/internal/app"
	"demo-service/logger"
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
		log.Fatal("error pars config")
	}

	logConf := logger.Config{
		LogLvl: cnf.Log.LogLevel,
	}

	logger := logger.NewLoger(logConf)

	app, err := app.NewApp(cnf, logger)
	if err != nil {
		logger.Error("error init app", "error", err)
		return
	}

	// для системных сигналов
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	servChError := make(chan error, 1)

	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			servChError <- err
		}
		close(servChError)
	}()

	select {
	case err := <-servChError:
		logger.Error("server crashed", "error", err)
		return
	case <-sigCh:
		logger.Info("Shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := app.Shutdown(ctx); err != nil {
			logger.Warn("Shutdown failed", "error", err)
			return
		}

		logger.Info("App shutdown complete")
	}

}
