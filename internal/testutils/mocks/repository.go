package mocks

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCategory(ctx context.Context, id int) (*model.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Category), args.Error(1)
}
func (m *MockRepository) GetCategoriesAll(ctx context.Context) (*model.AllCategories, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.AllCategories), args.Error(1)
}
func (m *MockRepository) CreateCategory(ctx context.Context, category *model.CreateCategoryRequest) (int, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(int), args.Error(1)
}
