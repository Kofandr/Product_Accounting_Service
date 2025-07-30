package model

type AllCategories struct {
	Categories []Category `json:"Categories"`
}

type Category struct {
	Id          int64  `json:"id" example:"4"`
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
