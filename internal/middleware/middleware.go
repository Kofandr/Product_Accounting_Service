package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"time"
)

type CtxLoggerKey struct {
}

func RequestLogger(logg *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := uuid.NewString()
			reqLog := logg.With("request_id", reqID)

			ctx := c.Request().Context()

			ctx = context.WithValue(ctx, CtxLoggerKey{}, reqLog)

			c.SetRequest(c.Request().WithContext(ctx))

			start := time.Now()

			err := next(c)

			duration := time.Since(start)

			req := c.Request()
			res := c.Response()

			logFields := []any{
				"method", req.Method, // HTTP-метод (GET, POST и т.д.)
				"path", req.URL.Path, // Путь URL (/sum, /multiply)
				"ip", c.RealIP(),
				"status", res.Status, // HTTP-статус ответа (200, 400 и т.д.)
				"duration", duration, // Время выполнения запроса
			}

			if err != nil {
				logFields = append(logFields, "err", err.Error())
				logg.Error("request failed", logFields...)
			} else {
				logg.Info("request handled", logFields...)
			}

			return err
		}
	}
}
