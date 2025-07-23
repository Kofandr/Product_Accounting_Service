package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) GetCategoriesAll(c echo.Context) error {

	return nil
}

func (handler *Handler) GetCategoryById(c echo.Context) error {
	logg := logger.MustLoggerFromCtx(c.Request().Context())

	ctx := c.Request().Context()

	categoryId := c.Param("id")

	var categories model.Categories
	err := handler.db.QueryRow(
		ctx,
		"SELECT id, name, description FROM categories WHERE id = $1",
		categoryId,
	).Scan(
		&categories.Id,
		&categories.Name,
		&categories.Description,
	)
	if err != nil {
		logg.Error("An error occurred while accessing the database", "err", err)
		return err
	}

	return c.JSON(http.StatusOK, categories)
}
