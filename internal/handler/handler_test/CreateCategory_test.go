package handler_test

import (
	"bytes"

	"github.com/Kofandr/Product_Accounting_Service/internal/apperrors"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"

	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		inputJSON      string
		mockReturn     int
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request",
			inputJSON:      `{"name": "Name", "description": "Description"}`,
			mockReturn:     3,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Id category": 3}`,
		},
		{
			name:           "Empty Name",
			inputJSON:      `{"name": "", "description": "Invalid"}`,
			mockReturn:     0,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"err": "Invalid JSON format"}`,
		},
		{
			name:           "Database Error",
			inputJSON:      `{"name": "Books", "description": "description"}`,
			mockReturn:     0,
			mockError:      apperrors.ErrConnectionFailed,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"err": "Server error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.Repository)

			c := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.inputJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			mockRepo.On("CreateCategory", mock.Anything, mock.Anything).Return(test.mockReturn, test.mockError)

			handler := handler.New(mockRepo)
			if err := handler.CreateCategory(cT); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}
}
