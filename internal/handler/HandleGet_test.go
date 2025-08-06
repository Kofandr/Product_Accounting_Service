package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerGet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		method       func(*Handler, echo.Context) error
		param        string
		mockMethod   string
		mockOn       int
		mockModel    interface{}
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		{
			name:       "GetCategoryById_Success",
			method:     (*Handler).GetCategoryById,
			param:      "1",
			mockOn:     1,
			mockMethod: "GetCategory",
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
			name:       "GetProduct_Success",
			method:     (*Handler).GetProduct,
			param:      "1",
			mockMethod: "GetProduct",
			mockOn:     1,
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
			name:       "GetProductsCategory_Success",
			method:     (*Handler).GetProductsCategory,
			param:      "1",
			mockMethod: "GetProductsCategory",
			mockOn:     1,
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
						"id":          1,
						"name":        "Bolls",
						"amount": 1,
						"category_id": 1
                    },
                    {
                        "id":          2,
						"name":        "R",
						"amount": 1,
						"category_id": 1
                    }
                ]
            }`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/"+test.param, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.param)

			mockBD := new(mocks.Repository)
			mockBD.On(test.mockMethod, mock.Anything, test.mockOn).Return(test.mockModel, test.mockReturn)
			handler := New(mockBD)
			if err := test.method(handler, c); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			assert.Equal(t, test.expectedCode, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}

}
