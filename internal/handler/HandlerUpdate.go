package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Kofandr/Product_Accounting_Service/internal/model"

	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func HandlerUpdate[T any](
	c echo.Context,
	updateFunc func(context.Context, int, *T) error,
	entityName string,
) error {
	logg := appctx.LoggerFromContext(c.Request().Context())
	ctx := c.Request().Context()

	id, err := parseIDParam(c)
	if err != nil {
		logg.Info("Invalid id",
			"entity", entityName,
			"error", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid id"})
	}

	var requestData T
	if err := c.Bind(&requestData); err != nil {
		logg.Error("Invalid JSON received",
			"entity", entityName,
			"error", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid JSON format"})
	}

	if err := c.Validate(requestData); err != nil {
		logg.Error("Validation failed",
			"entity", entityName,
			"error", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid request data"})
	}

	err = updateFunc(ctx, id, &requestData)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logg.Error("Entity not found",
				"entity", entityName,
				"id", id,
				"error", err)

			return c.JSON(http.StatusNotFound, map[string]string{"err": "Not found"})
		}

		logg.Error("Database error",
			"entity", entityName,
			"operation", "update",
			"id", id,
			"error", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"err": "Server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Request Status": "Changes completed",
	})
}

func (handler *Handler) UpdateCategory(c echo.Context) error {
	return HandlerUpdate[model.UpdateCategoryRequest](
		c,
		handler.db.UpdateCategory,
		"Category",
	)
}

func (handler *Handler) UpdateProduct(c echo.Context) error {
	return HandlerUpdate[model.UpdateProductRequest](
		c,
		handler.db.UpdateProduct,
		"Product",
	)
}
