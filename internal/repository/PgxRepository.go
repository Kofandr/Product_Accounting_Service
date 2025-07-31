package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicate = errors.New("duplicate entry")
)

type PgxRepository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *PgxRepository {
	return &PgxRepository{
		db,
	}
}

func (pgxRepository *PgxRepository) GetCategory(ctx context.Context, id int) (*model.Category, error) {
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

func (pgxRepository *PgxRepository) CreateCategory(ctx context.Context, category *model.CreateCategoryRequest) (int, error) {
	var id int
	err := pgxRepository.db.QueryRow(ctx,
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		category.Name,
		category.Description,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			if pgErr.Code == "23505" {
				return id, fmt.Errorf("%w: category '%s' already exists", ErrDuplicate, category.Name)
			}
		}
		return id, fmt.Errorf("failed to create category: %w", err)
	}
	return id, err
}

func (pgxRepository *PgxRepository) UpdateCategory(ctx context.Context, id int, update *model.UpdateCategoryRequest) error {
	query := `
        UPDATE categories
        SET
            name = COALESCE(NULLIF($1, ''), name),
            description = COALESCE(NULLIF($2, ''), description)
        WHERE id = $3
    `

	var namePtr, descPtr interface{}
	if update.Name != nil {
		namePtr = *update.Name
	} else {
		namePtr = nil
	}

	if update.Description != nil {
		descPtr = *update.Description
	} else {
		descPtr = nil
	}

	result, err := pgxRepository.db.Exec(ctx, query, namePtr, descPtr, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrUnsupported
	}

	return nil
}

func (pgxRepository *PgxRepository) DeleteCategory(ctx context.Context, id int) error {
	result, err := pgxRepository.db.Exec(ctx,
		"DELETE FROM  categories WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.ErrUnsupported
	}

	return nil
}
