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

func (handler *Handler) CreateCategory(c echo.Context) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	ctx := c.Request().Context()

	var category model.CreateCategoryRequest
	if err := c.Bind(&category); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}

		logg.Error("Invalid JSON received", "err", err)

		return c.JSON(http.StatusBadRequest, errResp)
	}

	if err := c.Validate(category); err != nil {
		errResp := map[string]string{"err": "Invalid request data"}

		logg.Error("Validation failed", "err", err)

		return c.JSON(http.StatusBadRequest, errResp)
	}

	id, err := handler.db.CreateCategory(ctx, &category)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			errResp := map[string]string{"err": fmt.Sprintf("Category with name '%s' already exists", category.Name)}

			logg.Warn("Duplicate category attempt", "name", category.Name, "err", err)

			return c.JSON(http.StatusConflict, errResp) // 409 Conflict
		}

		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", "err", err)

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"Id category": id,
	})
}
