package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/labstack/echo/v4"
)

func (handler *Handler) DeleteCategory(c echo.Context) error {
	logg := logger.MustLoggerFromCtx(c.Request().Context())

	ctx := c.Request().Context()
	return nil
}
