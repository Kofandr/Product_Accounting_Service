package handler

import "github.com/jackc/pgx/v5"

type Handler struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Handler {
	return &Handler{
		db,
	}
}
