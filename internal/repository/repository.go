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
}

type ProductRepository interface {
}
