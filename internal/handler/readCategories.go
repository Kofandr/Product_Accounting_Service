package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) GetCategoriesAll(c echo.Context) error {

	return nil
}

func (handler *Handler) GetCategoryByName(c echo.Context) error {
	ctx := c.Request().Context()
	row := handler.db.QueryRow(
		ctx,
		"SELECT name, description FROM categories WHERE name = $1",
		"boots",
	)

	var categories model.Categories
	err := row.Scan(
		&categories.Name,
		&categories.Description,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, categories)
}
