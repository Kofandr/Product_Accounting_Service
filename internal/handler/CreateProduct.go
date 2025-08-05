package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) CreateProduct(c echo.Context) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	ctx := c.Request().Context()

	var product model.CreateProductRequest
	if err := c.Bind(&product); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received", "err", err)
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if err := c.Validate(product); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received", "err", err)
		return c.JSON(http.StatusBadRequest, errResp)
	}

	be, err := handler.db.CategoryExists(ctx, product.CategoryId)
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

	id, err := handler.db.CreateProduct(ctx, &product)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}
		logg.Error("An error occurred while accessing the database", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Id product": id,
	})
}
