package handler_test

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/appvalidator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type setCaseGet struct {
	name         string
	param        string
	mockID       int
	mockModel    interface{}
	mockReturn   error
	expectedCode int
	expectedBody string
}

func runGetTest(t *testing.T, test setCaseGet, methodName string, methodFunc handlerFunc) {
	t.Helper()

	c := echo.New()

	c.Validator = &appvalidator.CustomValidator{Validator: validator.New()}

	req := httptest.NewRequest(http.MethodGet, "/"+test.param, nil)
	rec := httptest.NewRecorder()
	echoCtx := c.NewContext(req, rec)
	echoCtx.SetParamNames("id")
	echoCtx.SetParamValues(test.param)

	mockDB := new(mocks.Repository)
	mockDB.On(methodName, mock.Anything, test.mockID).Return(test.mockModel, test.mockReturn)

	handler := handler.New(mockDB)
	err := methodFunc(handler, echoCtx)

	require.NoError(t, err)

	assert.Equal(t, test.expectedCode, rec.Code)
	assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))
}

func TestGetCategoryByID(t *testing.T) {
	t.Parallel()

	tests := []setCaseGet{
		{
			name:   "Success",
			param:  "1",
			mockID: 1,
			mockModel: &model.Category{
				ID:          1,
				Name:        "Bolls",
				Description: "Bolls Description",
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": 1, "name": "Bolls", "description": "Bolls Description"}`,
		},
		{
			name:         "NotFound",
			param:        "999",
			mockID:       999,
			mockModel:    nil,
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Category not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			runGetTest(t, test, "GetCategory", (*handler.Handler).GetCategoryByID)
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	tests := []setCaseGet{
		{
			name:   "Success",
			param:  "1",
			mockID: 1,
			mockModel: &model.Product{
				ID:         1,
				Name:       "Bolls",
				Amount:     1,
				CategoryID: 1,
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": 1, "name": "Bolls", "amount": 1, "category_id": 1}`,
		},
		{
			name:         "NotFound",
			param:        "999",
			mockID:       999,
			mockModel:    nil,
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Product not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			runGetTest(t, test, "GetProduct", (*handler.Handler).GetProduct)
		})
	}
}

func TestGetProductsCategory(t *testing.T) {
	t.Parallel()

	tests := []setCaseGet{
		{
			name:   "Success",
			param:  "1",
			mockID: 1,
			mockModel: &model.ProductsCategory{
				Category: "Bolls",
				Products: []model.Product{
					{
						ID:         1,
						Name:       "Bolls",
						Amount:     1,
						CategoryID: 1,
					},
					{
						ID:         2,
						Name:       "R",
						Amount:     1,
						CategoryID: 1,
					},
				},
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{
				"Category": "Bolls",	
				"Products": [
					{
						"id": 1,
						"name": "Bolls",
						"amount": 1,
						"category_id": 1
					},
					{
						"id": 2,
						"name": "R",
						"amount": 1,
						"category_id": 1
					}
				]
			}`,
		},
		{
			name:         "NotFound",
			param:        "999",
			mockID:       999,
			mockModel:    nil,
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Category not found or empty"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runGetTest(t, test, "GetProductsCategory", (*handler.Handler).GetProductsCategory)
		})
	}
}
