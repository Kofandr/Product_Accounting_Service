package server

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/config"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/labstack/echo/v4"
	"log/slog"
	"strconv"
)

type Server struct {
	echo *echo.Echo
	addr string
	log  *slog.Logger
}

func New(log *slog.Logger, cfg *config.Configuration) *Server {
	serverEcho := echo.New()

	serverEcho.GET("/categories", handler.GetCategoriesAll)
	serverEcho.GET("/categories/:name", handler.GetCategoryByName)

	return &Server{serverEcho, (":" + strconv.Itoa(cfg.Port)), log}
}

func (server *Server) Start() error {
	server.log.Info("Starting server", "addr", server.addr)
	return server.echo.Start(server.addr)
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.log.Info("Shutting down server")
	return server.echo.Shutdown(ctx)
}
