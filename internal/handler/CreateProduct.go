package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Kofandr/Product_Accounting_Service/internal/repository"

	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/labstack/echo/v4"
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
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", "err", err)

		return c.JSON(http.StatusBadRequest, errResp)
	}

	Exist, err := handler.db.CategoryExists(ctx, product.CategoryID)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", "err", err)

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	if !Exist {
		errResp := map[string]string{"err": "Not found category"}

		logg.Error("Not found Category", "err", err)

		return c.JSON(http.StatusNotFound, errResp)
	}

	id, err := handler.db.CreateProduct(ctx, &product)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			errResp := map[string]string{"error": fmt.Sprintf("Category with name '%s' already exists", product.Name)}

			logg.Warn("Duplicate category attempt", "name", product.Name, "err", err)

			return c.JSON(http.StatusConflict, errResp) // 409 Conflict
		}

		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", "err", err)

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Id product": id,
	})
}
