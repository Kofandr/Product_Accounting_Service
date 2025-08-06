package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/handler/appValidator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }

func TestHandlerUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		method       func(*Handler, echo.Context) error
		param        string
		requestBody  string
		mockMethod   string
		mockID       int
		mockRequest  interface{}
		mockReturn   error
		expectedCode int
		expectedBody string
	}{
		// Тесты для UpdateCategory
		{
			name:        "UpdateCategory_Success",
			method:      (*Handler).UpdateCategory,
			param:       "1",
			requestBody: `{"name": "New Category", "description": "New Description"}`,
			mockMethod:  "UpdateCategory",
			mockID:      1,
			mockRequest: &model.UpdateCategoryRequest{
				Name:        stringPtr("New Category"),
				Description: stringPtr("New Description"),
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"Request Status":"Changes completed"}`,
		},
		{
			name:        "UpdateCategory_NotFound",
			method:      (*Handler).UpdateCategory,
			param:       "999",
			requestBody: `{"name": "Not Found Category"}`,
			mockMethod:  "UpdateCategory",
			mockID:      999,
			mockRequest: &model.UpdateCategoryRequest{
				Name: stringPtr("Not Found Category"),
			},
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Not found"}`,
		},
		{
			name:         "UpdateCategory_InvalidID",
			method:       (*Handler).UpdateCategory,
			param:        "invalid",
			requestBody:  `{"name": "Test"}`,
			mockMethod:   "",
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid id"}`,
		},
		{
			name:         "UpdateCategory_InvalidJSON",
			method:       (*Handler).UpdateCategory,
			param:        "1",
			requestBody:  `{"name": "Test", invalid}`,
			mockMethod:   "",
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid JSON format"}`,
		},

		{
			name:        "UpdateProduct_Success",
			method:      (*Handler).UpdateProduct,
			param:       "2",
			requestBody: `{"name": "Updated Product", "amount": 10, "category_id": 3}`,
			mockMethod:  "UpdateProduct",
			mockID:      2,
			mockRequest: &model.UpdateProductRequest{
				Name:       stringPtr("Updated Product"),
				Amount:     intPtr(10),
				CategoryID: intPtr(3),
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"Request Status":"Changes completed"}`,
		},
		{
			name:        "UpdateProduct_NotFound",
			method:      (*Handler).UpdateProduct,
			param:       "888",
			requestBody: `{"name": "Ghost Product"}`,
			mockMethod:  "UpdateProduct",
			mockID:      888,
			mockRequest: &model.UpdateProductRequest{
				Name: stringPtr("Ghost Product"),
			},
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Not found"}`,
		},
		{
			name:         "UpdateProduct_InvalidID",
			method:       (*Handler).UpdateProduct,
			param:        "nan",
			requestBody:  `{"name": "Test"}`,
			mockMethod:   "",
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid id"}`,
		},
		{
			name:         "UpdateProduct_InvalidJSON",
			method:       (*Handler).UpdateProduct,
			param:        "1",
			requestBody:  `{invalid json}`,
			mockMethod:   "",
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid JSON format"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(test.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.Validator = &appValidator.CustomValidator{Validator: validator.New()}
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.param)

			mockDB := new(mocks.Repository)

			// Настраиваем мок только если указан метод
			if test.mockMethod != "" {
				mockDB.On(
					test.mockMethod,
					mock.Anything,    // context.Context
					test.mockID,      // ID
					test.mockRequest, // request object
				).Return(test.mockReturn)
			}

			handler := New(mockDB)
			err := test.method(handler, c)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, rec.Code)
			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

			// Проверяем, что ожидаемые методы были вызваны
			if test.mockMethod != "" {
				mockDB.AssertExpectations(t)
			}
		})
	}
}
