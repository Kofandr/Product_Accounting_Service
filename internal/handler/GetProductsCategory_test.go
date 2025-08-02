package handler

import (
	"fmt"
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

func TestGetProductsCategory(t *testing.T) {
	tests := []struct {
		name                    string
		param                   string
		mockCategoryExistsOn    int
		mockCategoryExistsBool  bool
		mockCategoryExistsError error
		mockReturn              *model.ProductsCategory
		mockError               error
		expectedStatus          int
		expectedBody            string
	}{
		{
			name:                   "Valid Request",
			param:                  "1",
			mockCategoryExistsOn:   1,
			mockCategoryExistsBool: true,
			mockReturn: &model.ProductsCategory{
				Category: "Bolls",
				Products: []model.Product{
					{
						Id:         1,
						Name:       "Bolls",
						Amount:     1,
						CategoryId: 1,
					},
					{
						Id:         2,
						Name:       "R",
						Amount:     1,
						CategoryId: 1,
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
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
		{
			name:                   "Error BD",
			param:                  "1",
			mockCategoryExistsOn:   1,
			mockCategoryExistsBool: true,
			mockReturn:             &model.ProductsCategory{},
			mockError:              fmt.Errorf("erro"),
			expectedStatus:         http.StatusInternalServerError,
			expectedBody:           `{"err": "Server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mockBD := new(mocks.Repository)
			c := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			cT.SetParamNames("id")
			cT.SetParamValues(test.param)
			mockBD.On("CategoryExists", mock.Anything, test.mockCategoryExistsOn).Return(test.mockCategoryExistsBool, test.mockCategoryExistsError)
			mockBD.On("GetProductsCategory", mock.Anything, test.mockCategoryExistsOn).Return(test.mockReturn, test.mockError)
			handler := New(mockBD)
			handler.GetProductsCategory(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
