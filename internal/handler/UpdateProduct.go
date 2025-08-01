package handler

import (
	"database/sql"
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (handler Handler) UpdateProduct(c echo.Context) error {
	logg := logger.MustLoggerFromCtx(c.Request().Context())

	ctx := c.Request().Context()

	stringId := c.Param("id")

	productId, err := strconv.Atoi(stringId)
	if err != nil {
		errResp := map[string]string{"err": "Invalid id"}
		logg.Info("Invalid id", "err", err)
		return c.JSON(http.StatusBadRequest, errResp)
	}

	var product model.UpdateProductRequest
	if err := c.Bind(&product); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received")
		return c.JSON(http.StatusBadRequest, errResp)
	}

	err = handler.db.UpdateProduct(ctx, productId, &product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errResp := map[string]string{"err": "Not found"}
			logg.Error("Not found id", "err", err)
			return c.JSON(http.StatusNotFound, errResp)
		}

		errResp := map[string]string{"err": "Server error"}
		logg.Error("An error occurred while accessing the database", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Request Status": "Changes completed",
	})
}
