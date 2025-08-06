package handler

import (
	"context"
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/appctx"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteCategory(c echo.Context) error {
	return h.HandleDelete(c, h.db.DeleteCategory, "Category")
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	return h.HandleDelete(c, h.db.DeleteProduct, "Product")
}

func (h *Handler) HandleDelete(c echo.Context, deleteFunc func(context.Context, int) error, entity string) error {

	logg := appctx.LoggerFromContext(c.Request().Context())

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logg.Error("Invalid ID", "err", err)

		return c.JSON(http.StatusBadRequest, map[string]string{"err": "Invalid ID"})
	}

	err = deleteFunc(c.Request().Context(), id)
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

	return c.JSON(http.StatusOK, map[string]string{"message": entity + " deleted"})

}
