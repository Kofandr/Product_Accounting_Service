package server

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/config"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"log/slog"
	"strconv"
)

type Server struct {
	echo *echo.Echo
	addr string
	log  *slog.Logger
	db   *pgx.Conn
}

func New(log *slog.Logger, cfg *config.Configuration, db *pgx.Conn) *Server {
	serverEcho := echo.New()

	handler := handler.New(db)

	serverEcho.GET("/categories", handler.GetCategoriesAll)
	serverEcho.GET("/categories/:name", handler.GetCategoryByName)

	return &Server{serverEcho, (":" + strconv.Itoa(cfg.Port)), log, db}
}

func (server *Server) Start() error {
	server.log.Info("Starting server", "addr", server.addr)
	return server.echo.Start(server.addr)
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.log.Info("Shutting down server")
	return server.echo.Shutdown(ctx)
}
