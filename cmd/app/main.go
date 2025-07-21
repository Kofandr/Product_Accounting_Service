package main

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/config"
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/Kofandr/Product_Accounting_Service/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Mustload()
	mainLog := logger.New(cfg.LoggerLevel)

	mainServer := server.New(mainLog, cfg)

	go func() {
		if err := mainServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crash")

		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mainLog.Info("Shutting down...")
	if err := mainServer.Shutdown(ctx); err != nil {
		mainLog.Error("Shutdown failed", "error", err)
	} else {
		mainLog.Info("Server stopped")
	}

}
