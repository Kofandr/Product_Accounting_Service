package handler

import (
	"net/http"

	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/labstack/echo/v4"
)

func (handler *Handler) GetCategoriesAll(c echo.Context) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	ctx := c.Request().Context()

	categories, err := handler.db.GetCategoriesAll(ctx)
	if err != nil {
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", "err", err)

		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, categories)
}
