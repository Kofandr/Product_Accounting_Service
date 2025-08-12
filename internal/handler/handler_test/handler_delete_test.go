package handler_test

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/appvalidator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type setCaseDelete struct {
	name         string
	param        string
	mockID       int
	mockReturn   error
	expectedCode int
	expectedBody string
}

func runDeleteTest(t *testing.T, test setCaseDelete, methodName string, methodFunc handlerFunc) {
	t.Helper()

	c := echo.New()

	c.Validator = &appvalidator.CustomValidator{Validator: validator.New()}

	req := httptest.NewRequest(http.MethodDelete, "/"+test.param, nil)
	rec := httptest.NewRecorder()
	echoCtx := c.NewContext(req, rec)
	echoCtx.SetParamNames("id")
	echoCtx.SetParamValues(test.param)

	mockDB := new(mocks.Repository)

	if _, err := strconv.Atoi(test.param); err == nil {
		mockDB.On(methodName, mock.Anything, test.mockID).Return(test.mockReturn)
	}

	handler := handler.New(mockDB)
	err := methodFunc(handler, echoCtx)

	require.NoError(t, err)

	assert.Equal(t, test.expectedCode, rec.Code)
	assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

	if test.mockID > 0 {
		mockDB.AssertExpectations(t)
	}
}

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	tests := []setCaseDelete{
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
			runDeleteTest(t, test, "DeleteCategory", (*handler.Handler).DeleteCategory)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	t.Parallel()

	tests := []setCaseDelete{
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
			runDeleteTest(t, test, "DeleteProduct", (*handler.Handler).DeleteProduct)
		})
	}
}
