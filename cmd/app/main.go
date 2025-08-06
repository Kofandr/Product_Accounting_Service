package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kofandr/Product_Accounting_Service/internal/config"
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/Kofandr/Product_Accounting_Service/internal/server"
	"github.com/jackc/pgx/v5"
)

func main() {
	cfg := config.Mustload()
	logg := logger.New(cfg.LoggerLevel)

	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)

	}
	defer db.Close(context.Background())

	mainServer := server.New(logg, cfg, db)

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

	logg.Info("Shutting down...")
	if err := mainServer.Shutdown(ctx); err != nil {
		logg.Error("Shutdown failed", "error", err)
	} else {
		logg.Info("Server stopped")
	}

}
