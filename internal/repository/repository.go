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
	GetCategory(ctx context.Context, id int64) (*model.Category, error)
	GetCategoriesAll(ctx context.Context) (*model.AllCategories, error)
}

type ProductRepository interface {
}
