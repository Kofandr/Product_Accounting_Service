package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"

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

var errN error = sql.ErrNoRows

func TestGetCategoryById(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mokBD          mockDB
		expectedStatus int
		expectedBody   string
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
			expectedBody:   `{"id": 1, "name": "Bolls", "description": "Bolls Description"}`,
		},
		{
			name:  "Invalid id",
			param: "abc",
			mokBD: mockDB{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"err": "Invalid id"}`,
		},
		{
			name:  "not found",
			param: "999",
			mokBD: mockDB{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return &mockRow{
						ScanFunc: func(dest ...interface{}) error {
							return errN

						},
					}
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"err": "Not found"}`,
		},
		{
			name:  "Server error",
			param: "1",
			mokBD: mockDB{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return &mockRow{
						ScanFunc: func(dest ...interface{}) error {
							return fmt.Errorf("Server error")
						},
					}

				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"err": "Server error"}`,
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

			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))

		})
	}
}
