package handler_test

import (
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		param        string
		mockID       int
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success",
			param:        "1",
			mockID:       1,
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Category deleted"}`,
		},
		{
			name:         "NotFound",
			param:        "999",
			mockID:       999,
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Not found"}`,
		},
		{
			name:         "InvalidID",
			param:        "invalid",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid ID"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/categories/"+test.param, nil)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.param)

			mockDB := new(mocks.Repository)

			if _, err := strconv.Atoi(test.param); err == nil {
				mockDB.On("DeleteCategory", mock.Anything, test.mockID).Return(test.mockReturn)
			}

			h := handler.New(mockDB)
			err := h.DeleteCategory(echoCtx)

			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)
			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

			if test.mockID > 0 {
				mockDB.AssertExpectations(t)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		param        string
		mockID       int
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success",
			param:        "1",
			mockID:       1,
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Product deleted"}`,
		},
		{
			name:         "NotFound",
			param:        "999",
			mockID:       999,
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Not found"}`,
		},
		{
			name:         "InvalidID",
			param:        "invalid",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid ID"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/product/"+test.param, nil)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.param)

			mockDB := new(mocks.Repository)

			if _, err := strconv.Atoi(test.param); err == nil {
				mockDB.On("DeleteProduct", mock.Anything, test.mockID).Return(test.mockReturn)
			}

			h := handler.New(mockDB)
			err := h.DeleteProduct(echoCtx)

			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)
			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

			if test.mockID > 0 {
				mockDB.AssertExpectations(t)
			}
		})
	}
}
