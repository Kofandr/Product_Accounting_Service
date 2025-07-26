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
	Name        string `json:"name" `
	Description string `json:"description"`
}
