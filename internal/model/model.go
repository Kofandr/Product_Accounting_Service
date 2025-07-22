package model

type AllCategories struct{}

type Categories struct {
	Name        string `json:"name" example:"Name"`
	Description string `json:"description" example:"Description"`
}
