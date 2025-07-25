package mocks

import (
	"context"
	"github.com/Kofandr/Product_Accounting_Service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetCategory(ctx context.Context, id int64) (*model.Categories, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Categories), args.Error(1)
}
