package handler

import (
	"context"
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandlerGet[T any](
	c echo.Context,
	getFunc func(context.Context, int) (T, error),
	entity string,
	notFoundMsg string,
	invalidIDMsg string,
) error {
	logg := appctx.LoggerFromContext(c.Request().Context())
	ctx := c.Request().Context()

	id, err := parseIDParam(c)
	if err != nil {
		logg.Info("Invalid ID",
			"entity", entity,
			"input", c.Param("id"),
			"error", err.Error())

		return c.JSON(http.StatusBadRequest, map[string]string{"err": invalidIDMsg})
	}

	model, err := getFunc(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logg.Warn("Entity not found",
				"entity", entity,
				"id", id)

			return c.JSON(http.StatusNotFound, map[string]string{"err": notFoundMsg})
		}
		logg.Error("Database error",
			"operation", "Get"+entity,
			"id", id,
			"error", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]string{"err": "Server error"})
	}

	return c.JSON(http.StatusOK, model)
}

func (h *Handler) GetProduct(c echo.Context) error {
	return HandlerGet(
		c,
		h.db.GetProduct,
		"Product",
		"Product not found",
		"Invalid product ID",
	)
}

func (h *Handler) GetCategoryById(c echo.Context) error {
	return HandlerGet(
		c,
		h.db.GetCategory,
		"Category",
		"Category not found",
		"Invalid category ID",
	)
}

func (h *Handler) GetProductsCategory(c echo.Context) error {
	return HandlerGet(
		c,
		h.db.GetProductsCategory,
		"ProductsCategory",
		"Category not found or empty",
		"Invalid category ID",
	)
}
