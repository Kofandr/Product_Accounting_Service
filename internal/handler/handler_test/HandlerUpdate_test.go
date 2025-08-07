package handler_test

import (
	"strings"

	"github.com/Kofandr/Product_Accounting_Service/internal/appvalidator"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	// ... импорты ...
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type handlerFunc func(*handler.Handler, echo.Context) error
type setCaseUpdate struct {
	name         string
	param        string
	requestBody  string
	mockID       int
	mockRequest  interface{}
	mockReturn   error
	expectedCode int
	expectedBody string
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }

func runUpdateTest(t *testing.T, test setCaseUpdate, methodName string, methodFunc handlerFunc) {
	t.Helper()

	e := echo.New()

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(test.requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)

	echoCtx.SetParamNames("id")
	echoCtx.SetParamValues(test.param)

	e.Validator = &appvalidator.CustomValidator{Validator: validator.New()}

	mockDB := new(mocks.Repository)
	if test.mockRequest != nil {
		mockDB.On(methodName, mock.Anything, test.mockID, test.mockRequest).Return(test.mockReturn)
	}

	handler := handler.New(mockDB)
	err := methodFunc(handler, echoCtx)

	require.NoError(t, err)
	assert.Equal(t, test.expectedCode, rec.Code)
	assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

	if test.mockRequest != nil {
		mockDB.AssertExpectations(t)
	}
}

func TestUpdateCategory(t *testing.T) {
	t.Parallel()

	tests := []setCaseUpdate{
		{
			name:        "Success",
			param:       "1",
			requestBody: `{"name": "New Category", "description": "New Description"}`,
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
			name:        "NotFound",
			param:       "999",
			requestBody: `{"name": "Not Found Category"}`,
			mockID:      999,
			mockRequest: &model.UpdateCategoryRequest{
				Name: stringPtr("Not Found Category"),
			},
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Not found"}`,
		},
		{
			name:         "InvalidID",
			param:        "invalid",
			requestBody:  `{"name": "Test"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid id"}`,
		},
		{
			name:         "InvalidJSON",
			param:        "1",
			requestBody:  `{"name": "Test", invalid}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid JSON format"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			runUpdateTest(t, test, "UpdateCategory", (*handler.Handler).UpdateCategory)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	tests := []setCaseUpdate{
		{
			name:        "Success",
			param:       "2",
			requestBody: `{"name": "Updated Product", "amount": 10, "category_id": 3}`,
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
			name:        "NotFound",
			param:       "888",
			requestBody: `{"name": "Ghost Product"}`,
			mockID:      888,
			mockRequest: &model.UpdateProductRequest{
				Name: stringPtr("Ghost Product"),
			},
			mockReturn:   pgx.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"err":"Not found"}`,
		},
		{
			name:         "InvalidID",
			param:        "nan",
			requestBody:  `{"name": "Test"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid id"}`,
		},
		{
			name:         "InvalidJSON",
			param:        "1",
			requestBody:  `{invalid json}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"err":"Invalid JSON format"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			runUpdateTest(t, test, "UpdateProduct", (*handler.Handler).UpdateProduct)
		})
	}
}
