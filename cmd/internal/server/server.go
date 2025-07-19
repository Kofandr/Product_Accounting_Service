package server

import (
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Server struct {
	echo *echo.Echo
	addr string
	log  *slog.Logger
}

func New(log *slog.Logger) *Server {
	serverEcho := echo.New()

	return &Server{serverEcho, ":8080", log}
}
