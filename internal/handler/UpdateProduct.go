package handler

import (
	"errors"
	"net/http"

	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (handler Handler) UpdateProduct(c echo.Context) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	ctx := c.Request().Context()

	id, err := parseIDParam(c)
	if err != nil {
		logg.Info("Invalid id", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid id"})
	}

	var product model.UpdateProductRequest
	if err := c.Bind(&product); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received")
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if err := c.Validate(product); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received")
		return c.JSON(http.StatusBadRequest, errResp)
	}

	err = handler.db.UpdateProduct(ctx, id, &product)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
