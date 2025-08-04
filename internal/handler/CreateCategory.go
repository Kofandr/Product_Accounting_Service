package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/Kofandr/Product_Accounting_Service/internal/errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) CreateCategory(c echo.Context) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	ctx := c.Request().Context()

	var category model.CreateCategoryRequest
	if err := c.Bind(&category); err != nil {
		logg.Error("Invalid JSON received")
		resp := errors.MapErrorToResponse(err)
		return c.JSON(resp.Status, resp.Message)
	}
	if err := c.Validate(category); err != nil {
		logg.Error("Invalid JSON received")
		resp := errors.MapErrorToResponse(err)
		return c.JSON(resp.Status, resp.Message)
	}

	id, err := handler.db.CreateCategory(ctx, &category)
	if err != nil {
		logg.Error("An error occurred while accessing the database", "err", err)
		resp := errors.MapErrorToResponse(err)
		return c.JSON(resp.Status, resp.Message)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Id category": id,
	})
}
