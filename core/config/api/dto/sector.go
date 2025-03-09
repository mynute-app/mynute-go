package DTO

type Sector struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Your Company Sector Name"`
	Description string `json:"description" example:"The Company Sector Description"`
}
