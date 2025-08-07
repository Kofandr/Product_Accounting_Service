package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteHandlers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		method       func(*handler.Handler, echo.Context) error
		param        string
		mockMethod   string
		mockOn       int
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "DeleteCategory_Success",
			method:       (*handler.Handler).DeleteCategory,
			param:        "1",
			mockOn:       1,
			mockMethod:   "DeleteCategory",
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Category deleted"}`,
		},
		{
			name:         "DeleteProduct_Success",
			method:       (*handler.Handler).DeleteProduct,
			param:        "1",
			mockMethod:   "DeleteProduct",
			mockOn:       1,
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"Product deleted"}`,
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
			mockBD.On(test.mockMethod, mock.Anything, test.mockOn).Return(test.mockReturn)
			handler := handler.New(mockBD)

			if err := test.method(handler, c); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			assert.Equal(t, test.expectedCode, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}

}
