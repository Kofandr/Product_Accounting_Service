package handler

import (
	"context"
	"github.com/jackc/pgx/v5"
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
