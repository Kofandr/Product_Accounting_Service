package handler

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	db Repository
}

func New(db Repository) *Handler {
	return &Handler{
		db,
	}
}

type Repository interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Row interface {
	Scan(args ...any) error
}
