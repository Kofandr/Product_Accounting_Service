package repository

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/jackc/pgx/v5"
)

type PgxRepository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *PgxRepository {
	return &PgxRepository{
		db,
	}
}

func (pgxRepository *PgxRepository) GetCategory(ctx context.Context, id int64) (*model.Categories, error) {
	var categories model.Categories
	err := pgxRepository.db.QueryRow(ctx,
		"SELECT id, name, description FROM categories WHERE id = $1",
		id,
	).Scan(
		&categories.Id,
		&categories.Name,
		&categories.Description,
	)
	return &categories, err
}
