package handler

import (
	"fmt"
	"github.com/Kofandr/Product_Accounting_Service/internal/errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateCategory(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		inputJSON      string
		mockOn         int
		mockReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request",
			param:          "1",
			inputJSON:      `{"name": "Name", "description": "Description"}`,
			mockOn:         1,
			mockReturn:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Request Status": "Changes completed"}`,
		},
		{
			name:           "Invalid id",
			param:          "abc",
			mockOn:         1,
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"err": "Invalid id"}`,
		},
		{
			name:           "not found",
			param:          "999",
			mockOn:         999,
			mockReturn:     errors.ErrDBNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"err": "Not found"}`,
		},
		{
			name:           "Server error",
			param:          "1",
			mockOn:         1,
			mockReturn:     fmt.Errorf("erro"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"err": "Server error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockBD := new(mocks.Repository)
			c := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			cT.SetParamNames("id")
			cT.SetParamValues(test.param)
			mockBD.On("UpdateCategory", mock.Anything, test.mockOn, mock.Anything).Return(test.mockReturn, test.mockReturn)
			handler := New(mockBD)
			handler.UpdateCategory(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
