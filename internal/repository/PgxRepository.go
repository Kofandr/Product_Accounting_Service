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
		return pgx.ErrNoRows
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
		return pgx.ErrNoRows
	}

	return nil
}

func (pgxRepository *PgxRepository) GetProduct(ctx context.Context, id int) (*model.Product, error) {
	var product model.Product
	err := pgxRepository.db.QueryRow(ctx,
		"SELECT id, name, amount, categoryid FROM products WHERE id = $1",
		id,
	).Scan(
		&product.Id,
		&product.Name,
		&product.Amount,
		&product.CategoryId,
	)

	return &product, err
}

func (pgxRepository *PgxRepository) GetProductsCategory(ctx context.Context, categoryId int) (*model.ProductsCategory, error) {
	const query = `
        SELECT 
            c.name AS category_name,
            p.id, 
            p.name, 
            p.amount, 
            p.categoryid
        FROM categories c
        LEFT JOIN products p ON c.id = p.categoryid
        WHERE c.id = $1
    `

	var result model.ProductsCategory
	var products []model.Product

	rows, err := pgxRepository.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&result.Category,
			&p.Id,
			&p.Name,
			&p.Amount,
			&p.CategoryId,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if result.Category == "" {
		return nil, pgx.ErrNoRows
	}

	result.Products = products
	return &result, nil

}

func (pgxRepository *PgxRepository) CreateProduct(ctx context.Context, product *model.CreateProductRequest) (int, error) {
	var id int

	err := pgxRepository.db.QueryRow(ctx,
		"INSERT INTO products (name, amount, categoryid) VALUES ($1, $2, $3) RETURNING id",
		product.Name,
		product.Amount,
		product.CategoryId,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {

			if pgErr.Code == "23505" {
				return id, fmt.Errorf("%w: products  '%s' already exists", ErrDuplicate, product.Name)
			}
		}
		return id, fmt.Errorf("failed to create category: %w", err)
	}
	return id, err
}

func (pgxRepository *PgxRepository) UpdateProduct(ctx context.Context, id int, update *model.UpdateProductRequest) error {

	query := `
        UPDATE products
		SET
    		name = COALESCE(NULLIF($1, ''), name),
			amount = COALESCE(NULLIF($2, 0), amount),      -- 0 вместо ''
    		categoryId = COALESCE(NULLIF($3, 0), categoryId) -- 0 вместо ''
		WHERE id = $4
    `

	var namePtr, descPtr, catPtr interface{}
	if update.Name != nil {
		namePtr = *update.Name
	} else {
		namePtr = nil
	}

	if update.Amount != nil {
		descPtr = *update.Amount
	} else {
		descPtr = nil
	}

	if update.CategoryId != nil {
		catPtr = *update.CategoryId
	} else {
		catPtr = nil
	}

	result, err := pgxRepository.db.Exec(ctx, query, namePtr, descPtr, catPtr, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (pgxRepository *PgxRepository) DeleteProduct(ctx context.Context, id int) error {
	result, err := pgxRepository.db.Exec(ctx,
		"DELETE FROM products WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (pgxRepository *PgxRepository) CategoryExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`
	err := pgxRepository.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}
