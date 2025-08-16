//go:generate mockery --name=Repository --output=./mocks --case=underscore
package repository

import (
	"context"

	"github.com/Kofandr/Product_Accounting_Service/internal/model"
)

type Repository interface {
	CategoryRepository
	ProductRepository
}

type CategoryRepository interface {
	GetCategory(ctx context.Context, id int) (*model.Category, error)
	GetCategoriesAll(ctx context.Context) (*model.AllCategories, error)
	CreateCategory(ctx context.Context, category *model.CreateCategoryRequest) (int, error)
	UpdateCategory(ctx context.Context, id int, update *model.UpdateCategoryRequest) error
	DeleteCategory(ctx context.Context, id int) error
}

type ProductRepository interface {
	GetProduct(ctx context.Context, id int) (*model.Product, error)
	GetProductsCategory(ctx context.Context, categoryID int) (*model.ProductsCategory, error)
	CreateProduct(ctx context.Context, product *model.CreateProductRequest) (int, error)
	UpdateProduct(ctx context.Context, id int, update *model.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
	CategoryExists(ctx context.Context, id int) (bool, error)
}
