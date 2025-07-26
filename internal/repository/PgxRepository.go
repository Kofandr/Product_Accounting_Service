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

func (pgxRepository *PgxRepository) GetCategory(ctx context.Context, id int64) (*model.Category, error) {
	var categories model.Category
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

func (pgxRepository *PgxRepository) GetCategoriesAll(ctx context.Context) (*model.AllCategories, error) {

	rows, err := pgxRepository.db.Query(ctx, "SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category

	for rows.Next() {
		var c model.Category

		if err := rows.Scan(&c.Id, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &model.AllCategories{Categories: categories}, nil

}
