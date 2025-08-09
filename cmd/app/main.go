package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kofandr/Product_Accounting_Service/config"

	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/Kofandr/Product_Accounting_Service/internal/server"
)

func main() {
	cfg := config.MustLoad()
	logg := logger.New(cfg.LoggerLevel)

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	mainServer := server.New(logg, cfg, pool)

	go func() {
		if err := mainServer.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server crash")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShuttingDowntime)*time.Second)
	defer cancel()

	logg.Info("Shutting down...")

	if err := mainServer.Shutdown(ctx); err != nil {
		logg.Error("Shutdown failed", "error", err)
	} else {
		logg.Info("Server stopped")
	}
}
