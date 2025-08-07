package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/apperrors"
	"github.com/Kofandr/Product_Accounting_Service/internal/handler"

	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/Kofandr/Product_Accounting_Service/internal/repository/mocks"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerGetCategoriesAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		mockReturn     *model.AllCategories
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Request",
			mockReturn: &model.AllCategories{
				Categories: []model.Category{
					{
						ID:          1,
						Name:        "Bolls",
						Description: "Bolls Description",
					},
					{
						ID:          2,
						Name:        "R",
						Description: "R 00000000000000",
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: `{
                "Categories": [
                    {
                        "id": 1,
                        "name": "Bolls",
                        "description": "Bolls Description"
                    },
                    {
                        "id": 2,
                        "name": "R",
                        "description": "R 00000000000000"
                    }
                ]
            }`,
		},
		{
			name:           "Error BD",
			mockReturn:     &model.AllCategories{},
			mockError:      apperrors.ErrConnectionFailed,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"err": "Server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mockBD := new(mocks.Repository)
			c := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			mockBD.On("GetCategoriesAll", mock.Anything).Return(test.mockReturn, test.mockError)
			handler := handler.New(mockBD)
			if err := handler.GetCategoriesAll(cT); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			assert.Equal(t, test.expectedStatus, rec.Code)

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}

}
