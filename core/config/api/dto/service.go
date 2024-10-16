package DTO

type Service struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       int32      `json:"price"`
	Duration    int        `json:"duration"`
	Company     Company    `json:"company"`
	Branches    []Branch   `json:"branches"`
	Employees   []Employee `json:"employees"`
}
