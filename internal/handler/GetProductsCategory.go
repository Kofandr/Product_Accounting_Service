package handler

import (
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) GetProductsCategory(c echo.Context) error {
	logg := logger.MustLoggerFromCtx(c.Request().Context())
	ctx := c.Request().Context()

	id, err := parseIDParam(c)
	if err != nil {
		logg.Info("Invalid id", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid id"})
	}

	products, err := handler.db.GetProductsCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errResp := map[string]string{"err": "Not found category"}
			logg.Error("Not found category", "err", err)
			return c.JSON(http.StatusNotFound, errResp)
		}
		errResp := map[string]string{"err": "Server error"}
		logg.Error("Database error", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, products)

}
