package server

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/Kofandr/Product_Accounting_Service/internal/config"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler/appValidator"
	"github.com/Kofandr/Product_Accounting_Service/internal/middleware"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
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
	serverEcho.Validator = &appValidator.CustomValidator{Validator: validator.New()}

	serverEcho.GET("/categories", handler.GetCategoriesAll)
	serverEcho.GET("/categories/:id", handler.GetCategoryById)
	serverEcho.POST("/categories", handler.CreateCategory)
	serverEcho.PATCH("/categories/:id", handler.UpdateCategory)
	serverEcho.DELETE("/categories/:id", handler.DeleteCategory)

	serverEcho.GET("/product/:id", handler.GetProduct)
	serverEcho.GET("/products/:id", handler.GetProductsCategory)
	serverEcho.POST("/product", handler.CreateProduct)
	serverEcho.PATCH("/product/:id", handler.UpdateProduct)
	serverEcho.DELETE("/product/:id", handler.DeleteProduct)

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
