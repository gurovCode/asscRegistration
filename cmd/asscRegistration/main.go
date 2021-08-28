package main

import (
	"asscRegsitration/config"
	httpService "asscRegsitration/http"
	"asscRegsitration/model"
	"asscRegsitration/service"
	"context"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Read()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, _ := model.NewDB(logger, &cfg.Db)
	defer db.Stop()

	ctx := context.Background()

	svc := service.NewService(logger, cfg, db)

	httpSvc := httpService.NewService(ctx, logger, cfg, db, svc)
	httpClosedUnexpectedly := make(chan struct{})

	go func() {
		err := httpSvc.Start()

		if err == http.ErrServerClosed {
			logger.Info("http server was closed")
			return
		}

		logger.Error("http server closed unexpectedly", zap.Error(err))
		close(httpClosedUnexpectedly)
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logger.Info("Services started. Waiting for signal")

	select {
	case s := <-signals:
		logger.Info("Signal received", zap.String("signal", s.String()))
		httpSvc.Stop()
	case <-httpClosedUnexpectedly:
	}

}
