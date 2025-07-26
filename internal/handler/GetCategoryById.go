package handler

import (
	"database/sql"
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (handler *Handler) GetCategoryById(c echo.Context) error {
	logg := logger.MustLoggerFromCtx(c.Request().Context())

	ctx := c.Request().Context()

	stringId := c.Param("id")

	categoryId, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		errResp := map[string]string{"err": "Invalid id"}
		logg.Info("Invalid id", "err", err)
		return c.JSON(http.StatusBadRequest, errResp)
	}
	category, err := handler.db.GetCategory(ctx, categoryId)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			errResp := map[string]string{"err": "Not found"}
			logg.Error("Not found id", "err", err)
			return c.JSON(http.StatusNotFound, errResp)
		}
		errResp := map[string]string{"err": "Server error"}

		logg.Error("An error occurred while accessing the database", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, category)
}
