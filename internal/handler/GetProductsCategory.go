package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (handler *Handler) GetProductsCategory(c echo.Context) error {
	logg := logger.MustLoggerFromCtx(c.Request().Context())

	ctx := c.Request().Context()

	stringId := c.Param("id")

	categoryId, err := strconv.Atoi(stringId)
	if err != nil {
		errResp := map[string]string{"err": "Invalid id"}
		logg.Info("Invalid id", "err", err)

		return c.JSON(http.StatusBadRequest, errResp)
	}

	be, err := handler.db.CategoryExists(ctx, categoryId)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}
		logg.Error("An error occurred while accessing the database", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	if !be {
		errResp := map[string]string{"err": "Not found category"}
		logg.Error("Not found Category", "err", err)
		return c.JSON(http.StatusNotFound, errResp)
	}

	products, err := handler.db.GetProductsCategory(ctx, categoryId)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}
		logg.Error("An error occurred while accessing the database", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, products)

}
