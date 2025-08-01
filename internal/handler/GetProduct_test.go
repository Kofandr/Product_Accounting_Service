package handler

import (
	"fmt"
	"github.com/Kofandr/Product_Accounting_Service/internal/errors"
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

func TestGetProduct(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mockOn         int
		mockReturn     *model.Product
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Valid Request",
			param:  "1",
			mockOn: 1,
			mockReturn: &model.Product{
				Id:         1,
				Name:       "Bolls",
				Amount:     1,
				CategoryId: 1,
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 1, "name": "Bolls", "Amount": 1, "CategoryId": 1}`,
		},
		{
			name:   "Invalid id",
			param:  "abc",
			mockOn: 1,
			mockReturn: &model.Product{
				Id:         1,
				Name:       "Bolls",
				Amount:     1,
				CategoryId: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"err": "Invalid id"}`,
		},
		{
			name:   "not found",
			param:  "999",
			mockOn: 999,
			mockReturn: &model.Product{
				Id:         1,
				Name:       "Bolls",
				Amount:     1,
				CategoryId: 1,
			},
			mockErr:        errors.ErrDBNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"err": "Not found"}`,
		},
		{
			name:   "Server error",
			param:  "1",
			mockOn: 1,
			mockReturn: &model.Product{
				Id:         1,
				Name:       "Bolls",
				Amount:     1,
				CategoryId: 1,
			},
			mockErr:        fmt.Errorf("erro"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"err": "Server error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockBD := new(mocks.Repository)

			c := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			cT.SetParamNames("id")
			cT.SetParamValues(test.param)
			mockBD.On("GetProduct", mock.Anything, test.mockOn).Return(test.mockReturn, test.mockErr)
			handler := New(mockBD)
			handler.GetProduct(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}
}
