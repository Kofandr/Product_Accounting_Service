package handler

import (
	"fmt"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mockOn         int
		mockReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request",
			param:          "1",
			mockOn:         1,
			mockReturn:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Request Status": "Delete completed"}`,
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
			mockReturn:     pgx.ErrNoRows,
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
			t.Parallel()
			mockBD := new(mocks.Repository)
			c := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			cT.SetParamNames("id")
			cT.SetParamValues(test.param)
			mockBD.On("DeleteCategory", mock.Anything, test.mockOn).Return(test.mockReturn, test.mockReturn)
			handler := New(mockBD)
			handler.DeleteCategory(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
