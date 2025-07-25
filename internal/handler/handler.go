package handler

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/repository"
)

type Handler struct {
	db repository.Repository
}

func New(db repository.Repository) *Handler {
	return &Handler{
		db,
	}
}
