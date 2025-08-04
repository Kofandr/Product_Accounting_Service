package handler

import (
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (handler *Handler) UpdateCategory(c echo.Context) error {
	logg := appctx.LoggerFromContext(c.Request().Context())

	ctx := c.Request().Context()

	id, err := parseIDParam(c)
	if err != nil {
		logg.Info("Invalid id", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid id"})
	}

	var category model.UpdateCategoryRequest
	if err := c.Bind(&category); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received")
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if err := c.Validate(category); err != nil {
		errResp := map[string]string{"err": "Invalid JSON format"}
		logg.Error("Invalid JSON received")
		return c.JSON(http.StatusBadRequest, errResp)
	}

	err = handler.db.UpdateCategory(ctx, id, &category)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errResp := map[string]string{"err": "Not found"}
			logg.Error("Not found id", "err", err)
			return c.JSON(http.StatusNotFound, errResp)
		}

		errResp := map[string]string{"err": "Server error"}
		logg.Error("An error occurred while accessing the database", "err", err)
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Request Status": "Changes completed",
	})

}
