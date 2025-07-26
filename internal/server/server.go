package server

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/config"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/middleware"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"log/slog"
	"strconv"
)

type Server struct {
	echo *echo.Echo
	addr string
	logg *slog.Logger
	db   *pgx.Conn
}

func New(logg *slog.Logger, cfg *config.Configuration, db *pgx.Conn) *Server {
	serverEcho := echo.New()
	pgxRepository := repository.New(db)
	handler := handler.New(pgxRepository)
	serverEcho.Use(middleware.RequestLogger(logg))

	serverEcho.GET("/categories", handler.GetCategoriesAll)
	serverEcho.GET("/categories/:id", handler.GetCategoryById)
	serverEcho.POST("/categories", handler.CreateCategory)

	return &Server{serverEcho, (":" + strconv.Itoa(cfg.Port)), logg, db}
}

func (server *Server) Start() error {
	server.logg.Info("Starting server", "addr", server.addr)
	return server.echo.Start(server.addr)
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.logg.Info("Shutting down server")
	return server.echo.Shutdown(ctx)
}
