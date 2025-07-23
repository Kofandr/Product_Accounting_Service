package model

type AllCategories struct{}

type Categories struct {
	Id          int    `json:"id" example:"4"`
	Name        string `json:"name" example:"Name"`
	Description string `json:"description" example:"Description"`
}
