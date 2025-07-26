package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/testutils/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_CreateCategory(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		mockReturn     int
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Default",
			inputJSON:      `{"name": "Name", "description": "Description"}`,
			mockReturn:     3,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 3}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockBD := &mocks.MockRepository{}
			c := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.inputJSON))
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			mockBD.On("CreateCategory", mock.Anything).Return(test.mockReturn, test.mockError)
			handler := New(mockBD)
			handler.CreateCategory(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}
}
