package handler

import (
	"bytes"
	"errors"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name                    string
		inputJSON               string
		mockCategoryExistsOn    int
		mockCategoryExistsBool  bool
		mockCategoryExistsError error
		mockReturn              int
		mockError               error
		expectedStatus          int
		expectedBody            string
	}{
		{
			name:                    "Valid Request",
			inputJSON:               `{"name": "Name", "amount": 0, "category_id": 1}`,
			mockCategoryExistsOn:    1,
			mockCategoryExistsBool:  true,
			mockCategoryExistsError: nil,
			mockReturn:              3,
			mockError:               nil,
			expectedStatus:          http.StatusOK,
			expectedBody:            `{"Id product": 3}`,
		},
		{
			name:                    "Empty Name",
			inputJSON:               `{"name": "", "amount": , "category_id": }`,
			mockCategoryExistsOn:    1,
			mockCategoryExistsBool:  true,
			mockCategoryExistsError: nil,
			mockReturn:              0,
			mockError:               nil,
			expectedStatus:          http.StatusBadRequest,
			expectedBody:            `{"err": "Invalid JSON format"}`,
		},
		{
			name:                    "not found category",
			inputJSON:               `{"name": "Name", "amount": 0, "category_id": 999}`,
			mockCategoryExistsOn:    999,
			mockCategoryExistsBool:  false,
			mockCategoryExistsError: nil,
			mockReturn:              3,
			mockError:               nil,
			expectedStatus:          http.StatusNotFound,
			expectedBody:            `{"err": "Not found category"}`,
		},
		{
			name:                    "Database Error",
			inputJSON:               `{"name": "Name", "amount": 0, "category_id": 1}`,
			mockReturn:              0,
			mockCategoryExistsOn:    1,
			mockCategoryExistsBool:  true,
			mockCategoryExistsError: nil,
			mockError:               errors.New("connection failed"),
			expectedStatus:          http.StatusInternalServerError,
			expectedBody:            `{"err": "Server error"}`,
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
			mockRepo.On("CategoryExists", mock.Anything, test.mockCategoryExistsOn).Return(test.mockCategoryExistsBool, test.mockCategoryExistsError)
			mockRepo.On("CreateProduct", mock.Anything, mock.Anything).Return(test.mockReturn, test.mockError)
			handler := New(mockRepo)
			handler.CreateProduct(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}
}
