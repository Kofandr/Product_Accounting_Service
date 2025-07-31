package model

type AllCategories struct {
	Categories []Category `json:"Categories"`
}

type Category struct {
	Id          int    `json:"id" example:"4"`
	Name        string `json:"name" example:"Name"`
	Description string `json:"description" example:"Description"`
}
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type Product struct {
	Id         int    `json:"id" validate:"required,min=1"`
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Amount     int    `json:"Amount" validate:"required,min=0"`
	CategoryId int    `json:"CategoryId" validate:"required,min=1"`
}

type ProductsCategory struct {
	Category string    `json:"Category" validate:"required,min=2,max=100"`
	Products []Product `json:"Products"`
}

type CreateProductRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Amount     int    `json:"Amount" validate:"required,min=0"`
	CategoryId int    `json:"CategoryId" validate:"required,min=1"`
}

type UpdateProductRequest struct {
	Name       *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Amount     *int    `json:"Amount,omitempty" validate:"omitempty,min=0"`
	CategoryId *int    `json:"CategoryId,omitempty" validate:"omitempty,min=1"`
}
