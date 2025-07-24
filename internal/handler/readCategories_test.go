package handler

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"

	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/jackc/pgx/v5"

	"net/http"

	"testing"
)

// Реализуем мок DB с pointer receiver
type mockDB struct {
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (m *mockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return m.QueryRowFunc(ctx, sql, args...)
}

// Реализуем мок Row с pointer receiver
type mockRow struct {
	ScanFunc func(dest ...interface{}) error
}

func (m *mockRow) Scan(dest ...interface{}) error {
	return m.ScanFunc(dest...)
}

func TestGetCategoryById(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mokBD          mockDB
		expectedStatus int
		expectedBody   model.Categories
	}{
		{
			name:  "Normal",
			param: "1",
			mokBD: mockDB{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return &mockRow{
						ScanFunc: func(dest ...interface{}) error {
							*dest[0].(*int) = 1
							*dest[1].(*string) = "Bolls"
							*dest[2].(*string) = "Bolls Description"
							return nil
						},
					}
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: struct {
				Id          int    `json:"id" example:"4"`
				Name        string `json:"name" example:"Name"`
				Description string `json:"description" example:"Description"`
			}{Id: 1, Name: "Bolls", Description: "Bolls Description"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			handler := New(&test.mokBD)
			c := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			cT := c.NewContext(req, rec)
			cT.SetParamNames("id")
			cT.SetParamValues(test.param)

			handler.GetCategoryById(cT)

			assert.Equal(t, test.expectedStatus, rec.Code)

			var actualCategory model.Categories
			err := json.Unmarshal(rec.Body.Bytes(), &actualCategory)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			assert.Equal(t, test.expectedBody, actualCategory)

		})
	}
}
