package handler

import (
	"database/sql"
	"fmt"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/Kofandr/Product_Accounting_Service/internal/testutils/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"strings"

	"net/http"

	"testing"
)

var errN error = sql.ErrNoRows

func TestGetCategoryById(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mockOn         int64
		mockReturn     *model.Category
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Valid Request",
			param:  "1",
			mockOn: 1,
			mockReturn: &model.Category{
				Id:          1,
				Name:        "Bolls",
				Description: "Bolls Description",
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 1, "name": "Bolls", "description": "Bolls Description"}`,
		},
		{
			name:   "Invalid id",
			param:  "abc",
			mockOn: 1,
			mockReturn: &model.Category{
				Id:          1,
				Name:        "Bolls",
				Description: "Bolls Description",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"err": "Invalid id"}`,
		},
		{
			name:   "not found",
			param:  "999",
			mockOn: 999,
			mockReturn: &model.Category{
				Id:          1,
				Name:        "Bolls",
				Description: "Bolls Description",
			},
			mockErr:        errN,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"err": "Not found"}`,
		},
		{
			name:   "Server error",
			param:  "1",
			mockOn: 1,
			mockReturn: &model.Category{
				Id:          1,
				Name:        "Bolls",
				Description: "Bolls Description",
			},
			mockErr:        fmt.Errorf("erro"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"err": "Server error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockBD := &mocks.MockRepository{}

			c := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			cT.SetParamNames("id")
			cT.SetParamValues(test.param)
			mockBD.On("GetCategory", mock.Anything, test.mockOn).Return(test.mockReturn, test.mockErr)
			handler := New(mockBD)
			handler.GetCategoryById(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}
}
